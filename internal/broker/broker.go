package broker

import (
	"errors"
	"task2/internal/logger"
	"task2/internal/message"
	"task2/internal/subscriber"
)

var (
	// topic
	ErrTopicNotExist      = errors.New("topic does not exist")
	ErrNoStoredTopics     = errors.New("no topics stored")
	ErrNoTopicSubscribers = errors.New("no subscribers to topic")
	ErrTopicNameConflict  = errors.New("topic id already exists")
	ErrTopicIDEmpty       = errors.New("topic id cannot be empty")

	// subscriber
	ErrNoStoredSubscribers  = errors.New("no subscribers stored")
	ErrSubscriberNotFound   = errors.New("subscriber not found")
	ErrSubscriberNotActive  = errors.New("subscriber not active")
	ErrSubscriberIDEmpty    = errors.New("subscriber id cannot be empty")
	ErrSubAlreadySubscribed = errors.New("already subscribed to topic")
)

var (
	saveSubscriber         = make(chan *subscriber.Subscriber)
	getTopicAndSubscribers = make(chan getTopicAndSubscribersReq)
	getSubscriber          = make(chan getSubscriberReq)
	addSubscriberToTopic   = make(chan addSubscriberToTopicReq)
	createTopic            = make(chan string)
	publishMessage         = make(chan publishMessageReq)
)

type Broker struct {
	Subscribers subscriber.Subscribers
	Topics      map[string]subscriber.Subscribers
	logger      *logger.Logger
}

type getTopicAndSubscribersReq struct {
	topicID  string
	response chan map[string]subscriber.Subscribers
}

type getSubscriberReq struct {
	subID    string
	response chan *subscriber.Subscriber
}

type addSubscriberToTopicReq struct {
	sub     *subscriber.Subscriber
	topicID string
}

type publishMessageReq struct {
	message    *message.Message
	subscriber *subscriber.Subscriber
}

func NewBroker(logger *logger.Logger) *Broker {
	go func() {
		subscribers := subscriber.Subscribers{}
		topics := map[string]subscriber.Subscribers{}

		for {
			select {
			case sub := <-saveSubscriber:
				subscribers[sub.ID] = sub
			case req := <-getTopicAndSubscribers:
				subscribers, ok := topics[req.topicID]
				if !ok {
					req.response <- nil

					continue
				}
				req.response <- map[string]subscriber.Subscribers{req.topicID: subscribers}
			case req := <-getSubscriber:
				sub := subscribers[req.subID]
				req.response <- sub
			case req := <-addSubscriberToTopic:
				topics[req.topicID][req.sub.ID] = req.sub
			case topicID := <-createTopic:
				topics[topicID] = subscriber.Subscribers{}
			case req := <-publishMessage:
				req.subscriber.Messages <- req.message
			}
		}
	}()

	return &Broker{
		logger: logger,
	}
}

func (br *Broker) CreateSubscriber() *subscriber.Subscriber {
	br.logger.Log("Broker: CreateSubscriber")

	sub := subscriber.NewSubscriber(br.logger)

	saveSubscriber <- sub

	return sub
}

func (br *Broker) SubscribeToTopic(subID, topicID string) error {
	br.logger.Log("Broker: SubscribeToTopic")
	if topicID == "" {
		return ErrTopicIDEmpty
	}

	if subID == "" {
		return ErrSubscriberIDEmpty
	}

	response := make(chan map[string]subscriber.Subscribers)
	defer close(response)

	getTopicAndSubscribers <- getTopicAndSubscribersReq{
		topicID:  topicID,
		response: response,
	}

	topicAndSubscribers := <-response

	if len(topicAndSubscribers) == 0 {
		return ErrTopicNotExist
	}

	_, ok := topicAndSubscribers[topicID][subID]
	if ok {
		return ErrSubAlreadySubscribed
	}

	sub, err := br.GetSubscriber(subID)
	if err != nil {
		return err
	}

	if !sub.Active {
		return ErrSubscriberNotActive
	}

	addSubscriberToTopic <- addSubscriberToTopicReq{
		sub:     sub,
		topicID: topicID,
	}

	// need to update SubscribeToTopic to use similar monitor
	if err := sub.SubscribeToTopic(topicID, br.logger); err != nil {
		return err
	}

	return nil
}

func (br *Broker) CreateTopic(topicID string) error {
	br.logger.Log("Broker: CreateTopic")
	if topicID == "" {
		return ErrTopicIDEmpty
	}

	response := make(chan map[string]subscriber.Subscribers)
	defer close(response)

	getTopicAndSubscribers <- getTopicAndSubscribersReq{
		topicID:  topicID,
		response: response,
	}

	if <-response == nil {
		createTopic <- topicID

		return nil
	}

	return ErrTopicNameConflict
}

func (br *Broker) PublishMessage(topicID string, body interface{}) error {
	br.logger.Log("Broker: PublishMessage")
	if topicID == "" {
		return ErrTopicIDEmpty
	}
	response := make(chan map[string]subscriber.Subscribers)
	defer close(response)

	getTopicAndSubscribers <- getTopicAndSubscribersReq{
		topicID:  topicID,
		response: response,
	}

	topicAndSubscribers := <-response

	if len(topicAndSubscribers) == 0 {
		return ErrTopicNotExist
	}

	subscribers := topicAndSubscribers[topicID]

	if len(subscribers) == 0 {
		return ErrNoTopicSubscribers
	}

	msg, err := message.NewMessage(topicID, body, br.logger)
	if err != nil {
		return err
	}

	for _, s := range subscribers {
		publishMessage <- publishMessageReq{
			subscriber: s,
			message:    &msg,
		}
	}

	return nil
}

func (br *Broker) GetSubscriber(subID string) (*subscriber.Subscriber, error) {
	br.logger.Log("Broker: GetSubscriber")
	if subID == "" {
		return nil, ErrSubscriberIDEmpty
	}

	response := make(chan *subscriber.Subscriber)
	defer close(response)

	getSubscriber <- getSubscriberReq{
		subID:    subID,
		response: response,
	}

	sub := <-response

	if sub == nil {
		return nil, ErrSubscriberNotFound
	}

	return sub, nil
}
