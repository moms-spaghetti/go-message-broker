package protocols

import (
	"encoding/json"
	"errors"
	"net/http"
	"task2/internal/broker"
	"task2/internal/logger"
	"task2/internal/message"
	"task2/internal/subscriber"
)

func statusFromError(err error) int {
	switch {
	// ok
	case errors.Is(err, nil):
		return http.StatusOK
	// broker
	case errors.Is(err, broker.ErrTopicNotExist):
		return http.StatusNotFound
	case errors.Is(err, broker.ErrNoStoredSubscribers):
		return http.StatusNotFound
	case errors.Is(err, broker.ErrSubscriberNotFound):
		return http.StatusNotFound
	case errors.Is(err, broker.ErrSubscriberNotActive):
		return http.StatusUnauthorized
	case errors.Is(err, broker.ErrNoTopicSubscribers):
		return http.StatusNotFound
	case errors.Is(err, broker.ErrNoStoredTopics):
		return http.StatusNotFound
	case errors.Is(err, broker.ErrTopicNameConflict):
		return http.StatusConflict
	case errors.Is(err, broker.ErrTopicIDEmpty):
		return http.StatusBadRequest
	case errors.Is(err, broker.ErrSubscriberIDEmpty):
		return http.StatusBadRequest
	// messages
	case errors.Is(err, message.ErrMessageBodyEmpty):
		return http.StatusBadRequest
	case errors.Is(err, message.ErrMessageTopicEmpty):
		return http.StatusBadRequest
	// subscriber
	case errors.Is(err, subscriber.ErrSubNoTopics):
		return http.StatusNotFound
	case errors.Is(err, subscriber.ErrSubAlreadySubscribed):
		return http.StatusBadRequest
	case errors.Is(err, subscriber.ErrSubTopicNotExist):
		return http.StatusBadRequest
	case errors.Is(err, subscriber.ErrSubNoMessages):
		return http.StatusNotFound
	// server
	case errors.Is(err, ErrApiForbidden):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func BuildJsonResponse(err error, data interface{}, logger *logger.Logger) []byte {
	res := jsonResponse{
		Err:    "",
		Status: statusFromError(err),
		Data:   data,
	}

	if err != nil {
		logger.Log("BuildJsonResponse error:" + err.Error())
		res.Err = err.Error()
		res.Data = nil
	}

	out, err1 := json.Marshal(res)
	if err1 != nil {
		logger.Log("BuildJsonResponse error" + err.Error())
		res.Err = err1.Error()
		res.Status = http.StatusInternalServerError
	}

	logger.Log("BuildJsonResponse ok")

	return out
}
