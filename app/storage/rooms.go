package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Room struct {
	Name      string
	IsPrivate bool
	CreatedAt string
}

func EnsureRoom(name, pass string) error {
	if db == nil {
		return fmt.Errorf("db is not initialized")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("room name is empty")
	}

	_, _, err := GetRoom(name)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return CreateRoom(name, pass)
}

func CreateRoom(name, pass string) error {
	if db == nil {
		return fmt.Errorf("db is not initialized")
	}
	name = strings.TrimSpace(name)
	pass = strings.TrimSpace(pass)

	hash := ""
	if pass != "" {
		b, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
		hash = string(b)
	}

	_, err := db.Exec(
		"INSERT INTO rooms (name, salt, pass_hash, created_at) VALUES (?, ?, ?, ?)",
		name, "", hash, time.Now().UTC().Format(time.RFC3339),
	)
	return err
}

func GetRooms() ([]Room, error) {
	if db == nil {
		return nil, fmt.Errorf("db is not initialized")
	}

	rows, err := db.Query(`SELECT name, pass_hash, created_at FROM rooms ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Room
	for rows.Next() {
		var name, passHash, createdAt string
		if err := rows.Scan(&name, &passHash, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, Room{
			Name:      name,
			IsPrivate: strings.TrimSpace(passHash) != "",
			CreatedAt: createdAt,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func GetRoom(name string) (string, bool, error) {
	if db == nil {
		return "", false, fmt.Errorf("db is not initialized")
	}
	name = strings.TrimSpace(name)

	row := db.QueryRow(`SELECT name, pass_hash FROM rooms WHERE name = ?`, name)
	var n, passHash string
	if err := row.Scan(&n, &passHash); err != nil {
		return "", false, err
	}
	return n, strings.TrimSpace(passHash) != "", nil
}

func CheckRoomAccess(name, pass string) error {
	if db == nil {
		return fmt.Errorf("db is not initialized")
	}

	row := db.QueryRow(`SELECT pass_hash FROM rooms WHERE name = ?`, strings.TrimSpace(name))
	var passHash string
	if err := row.Scan(&passHash); err != nil {
		return err
	}

	passHash = strings.TrimSpace(passHash)
	if passHash == "" {
		return nil
	}

	pass = strings.TrimSpace(pass)
	if pass == "" {
		return fmt.Errorf("room is private: password required")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(pass)); err != nil {
		return fmt.Errorf("wrong password")
	}
	return nil
}
