package server

import (
	"fmt"

	"tcpChat/app/logger"
	"tcpChat/app/protocol"
)

func broadcastJSONToRoom(room string, exclude *Client, msg *protocol.Message) {
	line, err := protocol.Encode(msg)
	if err != nil {
		logger.L.Println("broadcast encode error:", err)
		return
	}

	recipients := snapshotRoomClients(room, exclude)

	for _, c := range recipients {
		if _, err := fmt.Fprint(c.Conn, line); err != nil {
			logger.L.Println("broadcast write error to", c.Nick, ":", err)
		}
	}
}

func broadcastJSON(from *Client, msg *protocol.Message) {
	broadcastJSONToRoom(from.Room, from, msg)
}
