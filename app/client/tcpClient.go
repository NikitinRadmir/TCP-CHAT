package client

import (
	"fmt"
	"net"
	"sync"

	"tcpChat/app/config"
	"tcpChat/app/logger"
)

func Start() error {
	cfg := config.Load()

	conn, err := net.Dial("tcp", cfg.ServerAddr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %w", cfg.ServerAddr, err)
	}
	logger.L.Println("Connected to", cfg.ServerAddr)

	var outMu sync.Mutex

	readErrCh := make(chan error, 1)

	go func() {
		readErrCh <- readLoop(conn, &outMu)
	}()

	printLocal(&outMu, "Type /help to see commands.")
	printPrompt(&outMu)

	writeErr := writeLoop(conn, &outMu)

	_ = conn.Close()
	_ = <-readErrCh

	logger.L.Println("Closing connection to server")

	return writeErr
}
