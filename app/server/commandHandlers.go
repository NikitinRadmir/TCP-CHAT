package server

import (
	"fmt"
	"strings"
	"time"

	"tcpChat/app/logger"
	"tcpChat/app/protocol"
	"tcpChat/app/storage"
)

func handleChatMessage(client *Client, msg *protocol.Message) {
	text := normalizeIncomingText(msg.Text)
	if text == "" {
		return
	}
	if len(text) > maxTextLen {
		text = text[:maxTextLen] + "…"
	}

	createdAt := time.Now().UTC().Format(time.RFC3339)

	if err := storage.SaveMessage(client.Nick, client.Room, text, createdAt); err != nil {
		logger.L.Println("failed to save message:", err)
	}

	out := &protocol.Message{
		Type:      "chat",
		Nick:      client.Nick,
		Room:      client.Room,
		Text:      text,
		CreatedAt: createdAt,
	}
	broadcastJSON(client, out)
}

func handleNickMessage(client *Client, msg *protocol.Message) {
	newNick := strings.TrimSpace(msg.NewNick)
	if newNick == "" {
		sendSystem(client, "nickname cannot be empty")
		return
	}
	if !nickRe.MatchString(newNick) {
		sendSystem(client, "invalid nickname. Use 3-20 chars: letters/digits/_/-, must start with letter/digit")
		return
	}

	oldNick := client.Nick
	if newNick == oldNick {
		sendSystem(client, fmt.Sprintf("you already use nickname %s", newNick))
		return
	}

	if err := setNick(client, newNick); err != nil {
		sendSystem(client, err.Error())
		return
	}

	sendSystem(client, fmt.Sprintf("your nickname is now %s", newNick))

	systemText := fmt.Sprintf("%s is now known as %s", oldNick, newNick)
	broadcastJSONToRoom(client.Room, client, &protocol.Message{
		Type:      "system",
		Room:      client.Room,
		Text:      systemText,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	})

	logger.L.Printf("Client changed nick: %s -> %s\n", oldNick, newNick)
}

func handleCreateRoomMessage(client *Client, msg *protocol.Message) {
	room := strings.TrimSpace(msg.Room)
	if !roomRe.MatchString(room) {
		sendSystem(client, "invalid room name. Use 3-32 chars: letters/digits/_/-, must start with letter/digit")
		return
	}

	pass := strings.TrimSpace(msg.Pass)

	if err := storage.CreateRoom(room, pass); err != nil {
		sendSystem(client, fmt.Sprintf("cannot create room '%s': %v", room, err))
		return
	}

	if pass != "" {
		sendSystem(client, fmt.Sprintf("room '%s' created (private)", room))
	} else {
		sendSystem(client, fmt.Sprintf("room '%s' created", room))
	}

	handleJoinMessage(client, &protocol.Message{Type: "join", Room: room, Pass: pass})
}

func handleJoinMessage(client *Client, msg *protocol.Message) {
	newRoom := strings.TrimSpace(msg.Room)
	if !roomRe.MatchString(newRoom) {
		sendSystem(client, "invalid room name. Use 3-32 chars: letters/digits/_/-, must start with letter/digit")
		return
	}

	if _, _, err := storage.GetRoom(newRoom); err != nil {
		sendSystem(client, fmt.Sprintf("room '%s' does not exist. Create it with: /create %s [PASS]", newRoom, newRoom))
		return
	}

	if err := storage.CheckRoomAccess(newRoom, msg.Pass); err != nil {
		sendSystem(client, fmt.Sprintf("cannot join '%s': %v", newRoom, err))
		return
	}

	if newRoom == client.Room {
		sendSystem(client, fmt.Sprintf("you are already in room %s", newRoom))
		return
	}

	oldRoom := setRoom(client, newRoom)

	broadcastJSONToRoom(oldRoom, client, &protocol.Message{
		Type:      "system",
		Room:      oldRoom,
		Text:      fmt.Sprintf("%s left the room", client.Nick),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	})

	broadcastJSONToRoom(newRoom, client, &protocol.Message{
		Type:      "system",
		Room:      newRoom,
		Text:      fmt.Sprintf("%s joined the room", client.Nick),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	})

	sendSystem(client, fmt.Sprintf("you joined room %s (previous: %s, users here: %d)", newRoom, oldRoom, clientCountInRoom(newRoom)))
	sendRoomHistory(client, newRoom, 30)
}

func handleLeaveMessage(client *Client) {
	if client.Room == "lobby" {
		sendSystem(client, "you are already in lobby")
		return
	}
	handleJoinMessage(client, &protocol.Message{Type: "join", Room: "lobby"})
}

func handleRoomsMessage(client *Client) {
	rooms, err := storage.GetRooms()
	if err != nil {
		sendSystem(client, fmt.Sprintf("failed to load rooms: %v", err))
		return
	}
	if len(rooms) == 0 {
		sendSystem(client, "no rooms")
		return
	}

	var b strings.Builder
	b.WriteString("rooms:\n")
	for _, r := range rooms {
		lock := ""
		if r.IsPrivate {
			lock = " 🔒"
		}
		b.WriteString(" - ")
		b.WriteString(r.Name)
		b.WriteString(lock)
		b.WriteString("\n")
	}
	sendSystem(client, strings.TrimSpace(b.String()))
}
