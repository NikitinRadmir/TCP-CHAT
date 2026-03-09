package main

import (
	"fmt"
	"os"

	"tcpChat/app/client"
	"tcpChat/app/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./app [server|client]")
		return
	}

	mode := os.Args[1]

	switch mode {
	case "server":
		if err := server.Start(); err != nil {
			fmt.Println("Server error:", err)
		}
	case "client":
		if err := client.Start(); err != nil {
			fmt.Println("Client error:", err)
		}
	default:
		fmt.Println("Unknown mode:", mode)
		fmt.Println("Usage: go run ./app [server|client]")
	}
}
