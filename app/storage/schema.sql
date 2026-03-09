PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS rooms (
                                     name       TEXT PRIMARY KEY,
                                     salt       TEXT NOT NULL DEFAULT '',
                                     pass_hash  TEXT NOT NULL DEFAULT '',
                                     created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
                                        id         INTEGER PRIMARY KEY AUTOINCREMENT,
                                        nick       TEXT NOT NULL,
                                        room       TEXT NOT NULL,
                                        text       TEXT NOT NULL,
                                        created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_messages_room_id ON messages(room, id);
