package protocol

import "encoding/json"

func Encode(msg *Message) (string, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func Decode(line string, msg *Message) error {
	return json.Unmarshal([]byte(line), msg)
}
