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
        if _, err := strconv.Atoi(os.Args[1]); err != nil {
            fmt.Println("[USAGE]: ./TCPChat $customport")
            return
        }
        port = os.Args[1]
    }

    server.Start(port)
}