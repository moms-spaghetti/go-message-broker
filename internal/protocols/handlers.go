package protocols

import (
	"encoding/json"
	"task2/internal/broker"
	"task2/internal/logger"
	"task2/internal/subscriber"
)

const (
	createSubscriber   = "createSubscriber"
	subscribeToTopic   = "subscribeToTopic"
	createTopic        = "createTopic"
	publishMessage     = "publishMessage"
	getMessages        = "getMessages"
	udpGetMessageCount = "udpGetMessageCount"
	udpGetNextMessage  = "udpGetNextMessage"
	udpCompleteMessage = "udpCompleteMessage"
)

type handlers struct {
	logger *logger.Logger
	br     *broker.Broker
}

func NewHandlers(logger *logger.Logger) *handlers {
	return &handlers{
		logger: logger,
		br:     broker.NewBroker(logger),
	}
}

func (h *handlers) createSubscriber(req jsonRequest) []byte {
	var (
		err error
		sub *subscriber.Subscriber
	)

	sub = h.br.CreateSubscriber()

	return BuildJsonResponse(err, sub, h.logger)
}

func (h *handlers) subscribeToTopic(req jsonRequest) []byte {
	var err error

	var query reqSubscribeToTopic
	err = json.Unmarshal(req.Query, &query)

	if err == nil {
		err = h.br.SubscribeToTopic(query.SubID, query.Topic)
	}

	return BuildJsonResponse(err, nil, h.logger)

}

func (h *handlers) createTopic(req jsonRequest) []byte {
	var err error

	var query reqAddTopic
	err = json.Unmarshal(req.Query, &query)

	if err == nil {
		err = h.br.CreateTopic(query.Name)
	}

	return BuildJsonResponse(err, nil, h.logger)
}

func (h *handlers) publishMessage(req jsonRequest) []byte {
	var err error

	var query reqMessageSubscribers
	err = json.Unmarshal(req.Query, &query)

	if err == nil {
		err = h.br.PublishMessage(query.Topic, query.Body)
	}

	return BuildJsonResponse(err, nil, h.logger)
}
