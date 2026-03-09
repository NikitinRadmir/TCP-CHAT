package storage

import (
	"database/sql"
	"embed"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

//go:embed schema.sql
var schemaFS embed.FS

func Init(dbPath string) error {
	var err error

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	sqlBytes, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("read embedded schema.sql: %w", err)
	}

	if err := execSchema(string(sqlBytes)); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}

	if err := EnsureRoom("lobby", ""); err != nil {
		return err
	}

	return nil
}

func Close() error {
	if db == nil {
		return nil
	}
	err := db.Close()
	db = nil
	return err
}
