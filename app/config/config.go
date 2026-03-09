package config

import "os"

type Config struct {
	ServerAddr string
	DBPath     string
}

func Load() Config {
	addr := os.Getenv("CHAT_ADDR")
	if addr == "" {
		addr = ":9000"
	}

	dbPath := os.Getenv("CHAT_DB_PATH")
	if dbPath == "" {
		dbPath = "chat.db"
	}

	return Config{
		ServerAddr: addr,
		DBPath:     dbPath,
	}
}
