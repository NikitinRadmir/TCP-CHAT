package client

import (
	"fmt"
	"net"

	"tcpChat/app/protocol"
)

func sendJSON(conn net.Conn, msg *protocol.Message) error {
	line, err := protocol.Encode(msg)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprint(conn, line); err != nil {
		return err
	}
	return nil
}
