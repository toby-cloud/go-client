package message

import (
	//	"fmt"
	//	"time"
	"encoding/json"
)

type Message struct {
	From    string            `json: "from,omitempty"`
	Payload map[string]string `json: "payload"`
	Tags    []string          `json: "tags"`
	Ack     string            `json: "ack"`
}

func NewMessage(from string, payload map[string]string, tags []string, ack string) *Message {
	m := &Message{
		From:    from,
		Payload: payload,
		Tags:    tags,
		Ack:     ack,
	}
	return m
}

func (m *Message) String() (string, error) {
	// omit the from field
	m.From = ""

	jsonString, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}
