package server

import (
    "fmt"
    "net"
    "sync"
)

const maxClients = 10

type Server struct {
    clients  map[*Client]bool
    history  []string
    mutex    sync.Mutex
    messages chan string
}

func NewServer() *Server {
    return &Server{
        clients:  make(map[*Client]bool),
        history:  make([]string, 0),
        messages: make(chan string, 100),
    }
}

func Start(port string) {
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        fmt.Printf("Error starting server: %v\n", err)
        return
    }
    defer listener.Close()

    fmt.Printf("Listening on the port :%s\n", port)

    server := NewServer()
    
    // Start message broadcaster
    go server.broadcaster()

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("Error accepting connection: %v\n", err)
            continue
        }

        server.mutex.Lock()
        if len(server.clients) >= maxClients {
            server.mutex.Unlock()
            conn.Write([]byte("Chat is full. Try again later.\n"))
            conn.Close()
            continue
        }
        server.mutex.Unlock()

        go server.handleConnection(conn)
    }
}

func (s *Server) handleConnection(conn net.Conn) {
    client := NewClient(conn, s)

    // Welcome message and get name
    if err := client.welcomeMessage(); err != nil {
        client.Close()
        return
    }

    // Check if name is already taken
    s.mutex.Lock()
    for existingClient := range s.clients {
        if existingClient.name == client.name {
            s.mutex.Unlock()
            conn.Write([]byte("Name already taken. Disconnecting.\n"))
            client.Close()
            return
        }
    }

    // Add client to server
    s.clients[client] = true
    s.mutex.Unlock()

    // Send chat history to new client
    s.sendHistory(client)

    // Notify everyone about new client
    joinMsg := formatSystemMessage(client.name + " has joined our chat.")
    s.messages <- joinMsg

    // Store join message in history
    s.mutex.Lock()
    s.history = append(s.history, joinMsg)
    s.mutex.Unlock()

    // Start client read/write routines
    go client.Write()
    go client.Read()
}

func (s *Server) broadcaster() {
    for msg := range s.messages {
        s.mutex.Lock()
        for client := range s.clients {
            select {
            case client.ch <- msg:
            default:
                // Skip if client channel is full
            }
        }
        s.mutex.Unlock()
    }
}

func (s *Server) broadcast(message string, sender *Client) {
    if message == "" {
        return
    }

    formattedMsg := formatMessage(sender.name, message)

    // Store in history
    s.mutex.Lock()
    s.history = append(s.history, formattedMsg)
    s.mutex.Unlock()

    // Broadcast to all clients except sender
    s.mutex.Lock()
    for client := range s.clients {
        if client != sender {
            select {
            case client.ch <- formattedMsg:
            default:
                // Skip if client channel is full
            }
        }
    }
    s.mutex.Unlock()
}

func (s *Server) removeClient(client *Client) {
    s.mutex.Lock()
    if _, exists := s.clients[client]; exists {
        delete(s.clients, client)
        s.mutex.Unlock()
        
        // Notify others
        leaveMsg := formatSystemMessage(client.name + " has left our chat.")
        s.messages <- leaveMsg
        
        // Store leave message in history
        s.mutex.Lock()
        s.history = append(s.history, leaveMsg)
        s.mutex.Unlock()
        
        client.Close()
    } else {
        s.mutex.Unlock()
    }
}

func (s *Server) sendHistory(client *Client) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    for _, msg := range s.history {
        client.ch <- msg
    }
}