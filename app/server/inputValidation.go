package server

import (
	"regexp"
	"strings"
)

var (
	nickRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{2,19}$`)
	roomRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{2,31}$`)
)

const (
	maxTextLen = 800
)

func normalizeIncomingText(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r < 32 || r == 127 {
			continue
		}
		out = append(out, r)
	}
	return strings.TrimSpace(string(out))
}
