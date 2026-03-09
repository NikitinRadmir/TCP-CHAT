package server

import (
	"fmt"
	"time"

	"tcpChat/app/logger"
	"tcpChat/app/protocol"
	"tcpChat/app/storage"
)

func sendSystem(to *Client, text string) {
	msg := &protocol.Message{
		Type:      "system",
		Text:      text,
		Room:      to.Room,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	line, err := protocol.Encode(msg)
	if err != nil {
		logger.L.Println("failed to marshal system message:", err)
		return
	}

	if _, err := fmt.Fprint(to.Conn, line); err != nil {
		logger.L.Println("failed to send system message to", to.Nick, ":", err)
	}
}

func sendRoomHistory(to *Client, room string, limit int) {
	history, err := storage.GetRoomHistory(room, limit)
	if err != nil {
		logger.L.Println("failed to load history:", err)
		return
	}

	if len(history) == 0 {
		sendSystem(to, fmt.Sprintf("room '%s' has no messages yet", room))
		return
	}

	sendSystem(to, fmt.Sprintf("history for room '%s' (last %d):", room, len(history)))

	for _, h := range history {
		msg := &protocol.Message{
			Type:      "chat",
			Nick:      h.Nick,
			Room:      h.Room,
			Text:      h.Text,
			CreatedAt: h.CreatedAt,
		}

		line, err := protocol.Encode(msg)
		if err != nil {
			logger.L.Println("failed to encode history msg:", err)
			continue
		}
		if _, err := fmt.Fprint(to.Conn, line); err != nil {
			logger.L.Println("failed to send history to", to.Nick, ":", err)
			return
		}
	}
}
