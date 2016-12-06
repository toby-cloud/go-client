package message

import (
	"encoding/json"
)

type Message struct {
	From    string                 `json:"from,omitempty"`
	Payload map[string]interface{} `json:"payload"`
	Tags    []string               `json:"tags"`
	Ack     string                 `json:"ack"`
	Id      string                 `json:"id,omitempty"`
	Sk      string                 `json:"sk,omitempty"`
}

func NewMessage(from string, payload map[string]interface{}, tags []string, ack string) *Message {
	m := &Message{
		From:    from,
		Payload: payload,
		Tags:    tags,
		Ack:     ack,
	}
	return m
}

func (m *Message) String() (string, error) {

	if len(m.Payload) == 0 {
		m.Payload = map[string]interface{}{}
	}
	if len(m.Tags) == 0 {
		m.Tags = []string{}
	}
	// omit the from field
	m.From = ""

	jsonString, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}
