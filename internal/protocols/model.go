package protocols

import (
	"encoding/json"
)

type jsonRequest struct {
	API     string          `json:"api"`
	Query   json.RawMessage `json:"query"`
	Payload json.RawMessage `json:"Payload"`
}

type jsonResponse struct {
	Err    string      `json:"err"`
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type reqSubscribeToTopic struct {
	SubID string `json:"subid"`
	Topic string `json:"topic"`
}

type reqAddTopic struct {
	Name string `json:"name"`
}

type reqMessageSubscribers struct {
	Topic string `json:"topic"`
	Body  string `json:"body"`
}

type reqMessagingQuery struct {
	Topic string `json:"topic"`
	SubID string `json:"subid"`
}
