package protocol

type Message struct {
	Type string `json:"type"`

	Nick      string `json:"nick,omitempty"`
	Text      string `json:"text,omitempty"`
	NewNick   string `json:"newNick,omitempty"`
	Room      string `json:"room,omitempty"`
	Pass      string `json:"pass,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}
