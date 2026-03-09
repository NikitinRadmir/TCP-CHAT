package server

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"tcpChat/app/logger"
)

type Client struct {
	Conn net.Conn
	Nick string
	Room string
}

var (
	clientsMu sync.RWMutex
	clients   = make(map[net.Conn]*Client)

	nextClientID atomic.Uint64
)

func addClient(conn net.Conn) *Client {
	id := nextClientID.Add(1)
	c := &Client{
		Conn: conn,
		Nick: fmt.Sprintf("user-%d", id),
		Room: "lobby",
	}

	clientsMu.Lock()
	clients[conn] = c
	total := len(clients)
	clientsMu.Unlock()

	logger.L.Println("Client added:", c.Nick, "room:", c.Room, "total:", total)
	return c
}

func removeClient(client *Client) {
	clientsMu.Lock()
	delete(clients, client.Conn)
	total := len(clients)
	clientsMu.Unlock()

	logger.L.Println("Client removed:", client.Nick, "room:", client.Room, "total:", total)
	_ = client.Conn.Close()
}

func setNick(client *Client, newNick string) error {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for _, c := range clients {
		if c == client {
			continue
		}
		if c.Nick == newNick {
			return fmt.Errorf("nickname '%s' is already taken", newNick)
		}
	}
	client.Nick = newNick
	return nil
}

func setRoom(client *Client, newRoom string) (oldRoom string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	oldRoom = client.Room
	client.Room = newRoom
	return oldRoom
}

func clientCountInRoom(room string) int {
	clientsMu.RLock()
	defer clientsMu.RUnlock()
	cnt := 0
	for _, c := range clients {
		if c.Room == room {
			cnt++
		}
	}
	return cnt
}

func snapshotRoomClients(room string, exclude *Client) []*Client {
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	out := make([]*Client, 0, 16)
	for _, c := range clients {
		if exclude != nil && c == exclude {
			continue
		}
		if c.Room != room {
			continue
		}
		out = append(out, c)
	}
	return out
}
