package storage

import (
	"fmt"
	"strings"
	"time"
)

type StoredMessage struct {
	Nick      string
	Room      string
	Text      string
	CreatedAt string
}

func SaveMessage(nick, room, text, createdAt string) error {
	if db == nil {
		return fmt.Errorf("db is not initialized")
	}
	nick = strings.TrimSpace(nick)
	room = strings.TrimSpace(room)
	text = strings.TrimSpace(text)
	createdAt = strings.TrimSpace(createdAt)
	if createdAt == "" {
		createdAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := db.Exec(
		"INSERT INTO messages (nick, room, text, created_at) VALUES (?, ?, ?, ?)",
		nick, room, text, createdAt,
	)
	return err
}

func GetRoomHistory(room string, limit int) ([]StoredMessage, error) {
	if db == nil {
		return nil, fmt.Errorf("db is not initialized")
	}
	if limit <= 0 {
		limit = 30
	}

	rows, err := db.Query(
		`SELECT nick, room, text, created_at
		 FROM messages
		 WHERE room = ?
		 ORDER BY id DESC
		 LIMIT ?`,
		room, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rev []StoredMessage
	for rows.Next() {
		var m StoredMessage
		if err := rows.Scan(&m.Nick, &m.Room, &m.Text, &m.CreatedAt); err != nil {
			return nil, err
		}
		rev = append(rev, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	return rev, nil
}
