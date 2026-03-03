package server

import (
    "bufio"
    "fmt"
    "net"
    "strings"
    "time"
)

type Client struct {
    conn     net.Conn
    name     string
    server   *Server
    ch       chan string
    scanner  *bufio.Scanner
    writer   *bufio.Writer
}

func NewClient(conn net.Conn, server *Server) *Client {
    return &Client{
        conn:    conn,
        server:  server,
        ch:      make(chan string, 100),
        scanner: bufio.NewScanner(conn),
        writer:  bufio.NewWriter(conn),
    }
}

func (c *Client) Read() {
    for c.scanner.Scan() {
        msg := strings.TrimSpace(c.scanner.Text())
        if msg == "" {
            continue
        }
        c.server.broadcast(msg, c)
    }

    // Client disconnected
    c.server.removeClient(c)
}

func (c *Client) Write() {
    for msg := range c.ch {
        c.writer.WriteString(msg)
        c.writer.Flush()
    }
}

func (c *Client) Close() {
    c.conn.Close()
    close(c.ch)
}

func (c *Client) welcomeMessage() error {
    logo := "Welcome to TCP-Chat!\n"
    logo += "         _nnnn_\n"
    logo += "        dGGGGMMb\n"
    logo += "       @p~qp~~qMb\n"
    logo += "       M|@||@) M|\n"
    logo += "       @,----.JM|\n"
    logo += "      JS^\\__/  qKL\n"
    logo += "     dZP        qKRb\n"
    logo += "    dZP          qKKb\n"
    logo += "   fZP            SMMb\n"
    logo += "   HZM            MMMM\n"
    logo += "   FqM            MMMM\n"
    logo += " __| \".        |\\dS\"qML\n"
    logo += " |    `.       | `' \\Zq\n"
    logo += "_)      \\.___,|     .'\n"
    logo += "\\____   )MMMMMP|   .'\n"
    logo += "     `-'       `--'\n"
    logo += "[ENTER YOUR NAME]: "

    _, err := c.writer.WriteString(logo)
    if err != nil {
        return err
    }
    c.writer.Flush()

    // Read client name
    if c.scanner.Scan() {
        name := strings.TrimSpace(c.scanner.Text())
        if name == "" {
            c.writer.WriteString("Name cannot be empty. Disconnecting.\n")
            c.writer.Flush()
            return fmt.Errorf("empty name")
        }
        c.name = name
    }
    return nil
}

func formatMessage(sender string, msg string) string {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    return fmt.Sprintf("[%s][%s]: %s\n", timestamp, sender, msg)
}

func formatSystemMessage(msg string) string {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    return fmt.Sprintf("[%s][System]: %s\n", timestamp, msg)
}