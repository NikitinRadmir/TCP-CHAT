package server

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"tcpChat/app/logger"
	"tcpChat/app/protocol"
)

func handleConnection(client *Client) {
	conn := client.Conn

	defer func() {
		room := client.Room
		nick := client.Nick

		removeClient(client)

		broadcastJSONToRoom(room, nil, &protocol.Message{
			Type:      "system",
			Room:      room,
			Text:      fmt.Sprintf("%s disconnected", nick),
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		})
		logger.L.Println("Closed connection from", nick)
	}()

	sendSystem(client, fmt.Sprintf("Welcome to the chat! Your nickname is: %s", client.Nick))
	sendSystem(client, helpText())
	sendRoomHistory(client, client.Room, 30)

	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0, 64*1024), 512*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var msg protocol.Message
		if err := protocol.Decode(line, &msg); err != nil {
			logger.L.Println("invalid json from", client.Nick, ":", err)
			sendSystem(client, "invalid message format, expected JSON")
			continue
		}

		dispatchMessage(client, &msg)
	}

	if err := scanner.Err(); err != nil {
		logger.L.Println("read error from", client.Nick, ":", err)
	}
}

func dispatchMessage(client *Client, msg *protocol.Message) {
	switch msg.Type {
	case "chat":
		handleChatMessage(client, msg)
	case "nick":
		handleNickMessage(client, msg)
	case "create":
		handleCreateRoomMessage(client, msg)
	case "join":
		handleJoinMessage(client, msg)
	case "leave":
		handleLeaveMessage(client)
	case "rooms":
		handleRoomsMessage(client)
	case "help":
		sendSystem(client, helpText())
	default:
		sendSystem(client, "unknown message type: "+msg.Type)
	}
}

func helpText() string {
	return "Commands: /help, /nick NAME, /rooms, /create ROOM [PASS], /join ROOM [PASS], /leave, /quit"
}
