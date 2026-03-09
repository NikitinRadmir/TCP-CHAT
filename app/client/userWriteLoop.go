package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"tcpChat/app/protocol"
)

func writeLoop(conn net.Conn, outMu *sync.Mutex) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("stdin read error: %w", err)
		}
		text = strings.TrimSpace(text)

		if text == "" {
			printPrompt(outMu)
			continue
		}

		if text == "/quit" || text == "/exit" {
			return nil
		}

		if text == "/help" {
			printHelp(outMu)
			printPrompt(outMu)
			continue
		}

		fields := strings.Fields(text)
		if len(fields) > 0 && strings.HasPrefix(fields[0], "/") {
			if err := handleSlashCommand(conn, outMu, fields); err != nil {
				return err
			}
			printPrompt(outMu)
			continue
		}

		msg := &protocol.Message{
			Type: "chat",
			Text: text,
		}
		if err := sendJSON(conn, msg); err != nil {
			return fmt.Errorf("write chat message: %w", err)
		}

		printPrompt(outMu)
	}
}

func handleSlashCommand(conn net.Conn, outMu *sync.Mutex, fields []string) error {
	switch fields[0] {
	case "/nick":
		if len(fields) < 2 {
			printLocal(outMu, "usage: /nick NAME")
			return nil
		}
		msg := &protocol.Message{Type: "nick", NewNick: fields[1]}
		if err := sendJSON(conn, msg); err != nil {
			return fmt.Errorf("write nick message: %w", err)
		}
		return nil

	case "/rooms":
		msg := &protocol.Message{Type: "rooms"}
		if err := sendJSON(conn, msg); err != nil {
			return fmt.Errorf("write rooms message: %w", err)
		}
		return nil

	case "/create":
		if len(fields) < 2 {
			printLocal(outMu, "usage: /create ROOM [PASS]")
			return nil
		}
		pass := ""
		if len(fields) >= 3 {
			pass = fields[2]
		}
		msg := &protocol.Message{Type: "create", Room: fields[1], Pass: pass}
		if err := sendJSON(conn, msg); err != nil {
			return fmt.Errorf("write create message: %w", err)
		}
		return nil

	case "/join":
		if len(fields) < 2 {
			printLocal(outMu, "usage: /join ROOM [PASS]")
			return nil
		}
		pass := ""
		if len(fields) >= 3 {
			pass = fields[2]
		}
		msg := &protocol.Message{Type: "join", Room: fields[1], Pass: pass}
		if err := sendJSON(conn, msg); err != nil {
			return fmt.Errorf("write join message: %w", err)
		}
		return nil

	case "/leave":
		msg := &protocol.Message{Type: "leave"}
		if err := sendJSON(conn, msg); err != nil {
			return fmt.Errorf("write leave message: %w", err)
		}
		return nil

	default:
		printLocal(outMu, "unknown command: "+fields[0]+" (type /help)")
		return nil
	}
}
