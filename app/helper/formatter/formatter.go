package formatter

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"tcpChat/app/protocol"
)

func Format(msg *protocol.Message) string {
	switch msg.Type {
	case "chat":
		return formatChat(msg)
	case "system":
		return formatSystem(msg)
	default:
		return fmt.Sprintf("%s [server/%s] %s", formatTime(msg.CreatedAt), renderToken(msg.Type), renderText(msg.Text))
	}
}

func formatChat(msg *protocol.Message) string {
	ts := formatTime(msg.CreatedAt)
	room := msg.Room
	if room == "" {
		room = "-"
	}
	nick := msg.Nick
	if nick == "" {
		nick = "unknown"
	}
	return fmt.Sprintf("%s #%s | %s: %s", ts, renderToken(room), renderToken(nick), renderText(msg.Text))
}

func formatSystem(msg *protocol.Message) string {
	ts := formatTime(msg.CreatedAt)
	text := renderText(msg.Text)
	if msg.Room != "" {
		return fmt.Sprintf("%s --- %s (#%s) ---", ts, text, renderToken(msg.Room))
	}
	return fmt.Sprintf("%s --- %s ---", ts, text)
}

func formatTime(createdAt string) string {
	createdAt = strings.TrimSpace(createdAt)
	if createdAt == "" {
		return "--:--"
	}

	if t, err := time.Parse(time.RFC3339Nano, createdAt); err == nil {
		return t.Local().Format("15:04")
	}
	if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
		return t.Local().Format("15:04")
	}

	if t, err := time.ParseInLocation("2006-01-02 15:04:05", createdAt, time.UTC); err == nil {
		return t.Local().Format("15:04")
	}

	if len(createdAt) >= 16 {
		chunk := createdAt[11:16]
		if strings.Count(chunk, ":") == 1 {
			return chunk
		}
	}
	return "--:--"
}

func renderToken(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "-"
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' || r == '\t' {
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	if out == "" {
		return "-"
	}
	if len(out) > 40 {
		return out[:40] + "…"
	}
	return out
}

func renderText(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			b.WriteRune(' ')
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		b.WriteRune(r)
	}
	out := strings.TrimSpace(b.String())
	if len(out) > 800 {
		return out[:800] + "…"
	}
	return out
}
