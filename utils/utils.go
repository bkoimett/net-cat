package main

import (
    "fmt"
    "os"
    "strconv"

    "net-cat/server"
)

func main() {
    port := "8989" // default port

    if len(os.Args) > 2 {
        fmt.Println("[USAGE]: ./TCPChat $port")
        return
    }

    if len(os.Args) == 2 {
        // Validate port number
        portNum, err := strconv.Atoi(os.Args[1])
        if err != nil || portNum < 1 || portNum > 65535 {
            fmt.Println("[USAGE]: ./TCPChat $port")
            fmt.Println("Port must be a number between 1 and 65535")
            return
        }
        port = os.Args[1]
    }

    server.Start(port)
}