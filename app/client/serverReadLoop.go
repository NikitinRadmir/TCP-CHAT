package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"

	"tcpChat/app/helper/formatter"
	"tcpChat/app/logger"
	"tcpChat/app/protocol"
)

func readLoop(conn net.Conn, outMu *sync.Mutex) error {
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0, 64*1024), 512*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var msg protocol.Message
		if err := protocol.Decode(line, &msg); err != nil {
			printLine(outMu, fmt.Sprintf("[server/raw] %s", line))
			continue
		}

		printLine(outMu, formatter.Format(&msg))
	}

	if err := scanner.Err(); err != nil {
		logger.L.Println("read error from server:", err)
		printLocal(outMu, "connection error. exiting...")
		return err
	}

	logger.L.Println("server closed the connection")
	printLocal(outMu, "server closed the connection. exiting...")
	return nil
}
