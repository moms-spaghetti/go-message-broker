package message

import (
	"crypto/rand"
	"errors"
	"fmt"
	"task2/internal/logger"
)

var (
	ErrMessageBodyEmpty  = errors.New("message body empty")
	ErrMessageTopicEmpty = errors.New("message topic empty")
)

type Message struct {
	ID    string      `json:"ID"`
	Topic string      `json:"topic"`
	Body  interface{} `json:"body"`
	Done  bool        `json:"done"`
}

func NewMessage(topic string, body interface{}, logger *logger.Logger) (Message, error) {
	logger.Log("creating new message")

	if body == "" {
		return Message{}, ErrMessageBodyEmpty
	}

	if topic == "" {
		return Message{}, ErrMessageTopicEmpty
	}

	m := Message{
		ID:    makeMessageID(),
		Topic: topic,
		Body:  body,
		Done:  false,
	}

	logger.Log("message created id: " + m.ID)

	return m, nil
}

func makeMessageID() string {
	id := make([]byte, 4)
	rand.Read(id)

	return fmt.Sprintf("%X", id[0:4])
}
