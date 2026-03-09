package server

import (
	"fmt"
	"net"

	"tcpChat/app/config"
	"tcpChat/app/logger"
	"tcpChat/app/storage"
)

func Start() error {
	cfg := config.Load()
	logger.L.Println("Starting server on", cfg.ServerAddr)

	if err := storage.Init(cfg.DBPath); err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}
	defer func() {
		_ = storage.Close()
	}()

	ln, err := net.Listen("tcp", cfg.ServerAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer ln.Close()

	logger.L.Println("Server is listening...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.L.Println("accept error:", err)
			continue
		}

		logger.L.Println("New connection from", conn.RemoteAddr())

		client := addClient(conn)
		go handleConnection(client)
	}
}
