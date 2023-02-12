package subscriber

import (
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"task2/internal/logger"
	"task2/internal/message"
)

var (
	ErrSubNoTopics          = errors.New("subscriber has no topics")
	ErrSubAlreadySubscribed = errors.New("already subscribed to topic")
	ErrSubTopicNotExist     = errors.New("topic does not exist")
	ErrSubNoMessages        = errors.New("no messages for topic")
)

type Subscriber struct {
	ID          string                        `json:"ID"`
	Topics      map[string][]*message.Message `json:"Topics"`
	Active      bool                          `json:"Active"`
	Messages    chan *message.Message         `json:"-"`
	done        chan struct{}
	topicsMutex *sync.RWMutex
}

type Subscribers map[string]*Subscriber

func NewSubscriber(logger *logger.Logger) *Subscriber {
	logger.Log("creating new subscriber")
	var sub Subscriber

	sub.ID = makeSubscriberID()
	sub.Topics = make(map[string][]*message.Message)
	sub.Messages = make(chan *message.Message)
	sub.done = make(chan struct{})
	sub.topicsMutex = &sync.RWMutex{}

	sub.Active = true

	go func(sub *Subscriber) {
		for {
			select {
			case <-sub.done:
				return
			default:
				newMessage := <-sub.Messages
				if newMessage != nil {
					sub.Topics[newMessage.Topic] = append(sub.Topics[newMessage.Topic], newMessage)
				}
			}
		}
	}(&sub)

	logger.Log("subscriber created id: " + sub.ID)

	return &sub

}

func (s *Subscriber) SubscribeToTopic(topicID string, logger *logger.Logger) error {
	logger.Log(fmt.Sprintf("subscriber %s subscribed to %s", s.ID, topicID))

	s.topicsMutex.Lock()
	defer s.topicsMutex.Unlock()

	_, ok := s.Topics[topicID]
	if !ok {
		s.Topics[topicID] = []*message.Message{}

		return nil
	}

	return ErrSubAlreadySubscribed
}

func (s *Subscriber) GetMessageCount(topicID string) int {
	s.topicsMutex.Lock()
	defer s.topicsMutex.Unlock()

	m := s.Topics[topicID]

	return len(m)
}

func (s *Subscriber) GetNextMessage(topicID string) message.Message {
	s.topicsMutex.Lock()
	defer s.topicsMutex.Unlock()

	return *s.Topics[topicID][0]
}

func (s *Subscriber) DeleteMessage(topicID, messageID string, logger *logger.Logger) error {
	if len(s.Topics[topicID]) == 0 {
		return ErrSubNoMessages
	}

	updated := make([]*message.Message, len(s.Topics[topicID])-1)

	i := 0
	for _, m := range s.Topics[topicID] {
		if m.ID == messageID {
			continue
		}

		updated[i] = m
		i++
	}

	logger.Log("sub messages updated")
	s.Topics[topicID] = updated

	return nil
}

func (s *Subscriber) ValidateTopicAndMessageCount(topicID string) error {
	messages, ok := s.Topics[topicID]
	if !ok {
		return ErrSubTopicNotExist
	}

	if len(messages) == 0 {
		return ErrSubNoMessages
	}

	return nil
}

func makeSubscriberID() string {
	id := make([]byte, 4)
	rand.Read(id)

	return fmt.Sprintf("%X", id[0:4])
}
