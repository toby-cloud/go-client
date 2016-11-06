package message

import (
	//	"fmt"
	//	"time"
	"encoding/json"
)

type Message struct {
	Message     string   `json:"message"`
	MessageType string   `json:"messageType"`
	Tags        []string `json:"tags"`
	AckTag      string   `json:"ackTag"`
}

func NewMessage(message, messageType string, tags []string, ackTag string) *Message {
	o := &Message{
		Message:     message,
		MessageType: messageType,
		Tags:        tags,
		AckTag:      ackTag,
	}
	return o
}

func (m *Message) String() (string, error) {
	jsonString, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}
