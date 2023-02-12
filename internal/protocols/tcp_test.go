package protocols

// import (
// 	"encoding/json"
// 	"net"
// 	"net/http"
// 	"reflect"
// 	"task2/internal/logger"
// 	"task2/internal/metrics"
// 	"task2/internal/subscriber"
// 	"testing"
// )

// func TestTCPServer_Start(t *testing.T) {
// 	tests := []struct {
// 		name string
// 	}{
// 		{
// 			name: "ok",
// 		},
// 	}
// 	for _, tt := range tests {
// 		logger := logger.NewLogger()
// 		metrics := metrics.NewMetrics(logger)
// 		handlers := NewHandlers(logger)
// 		tcp := NewTCP(logger, metrics, handlers)

// 		logger.StartNoopLogger()
// 		tcp.Start()
// 		metrics.Start()

// 		t.Run(tt.name, func(t *testing.T) {
// 			conn, err := net.Dial("tcp", ":8181")
// 			if err != nil {
// 				t.Errorf("failed to dial tcp server error: %v", err)
// 			}

// 			conn.Close()
// 		})
// 		tcp.Stop()
// 	}
// }

// type tcpTestRequest struct {
// 	Method  string                 `json:"method"`
// 	API     string                 `json:"aPI"`
// 	Query   map[string]interface{} `json:"query"`
// 	Payload map[string]interface{} `json:"payload"`
// }

// func TestTCPServer_TCPServer(t *testing.T) {
// 	type args struct {
// 		request tcpTestRequest
// 	}
// 	tests := []struct {
// 		name             string
// 		args             args
// 		addtopic         bool
// 		addSubscriber    bool
// 		subscribeTopic   bool
// 		mockSubscriberID bool
// 		want             jsonResponse
// 	}{
// 		{
// 			name: "ok - add topic",
// 			args: args{
// 				request: tcpTestRequest{
// 					Method: "POST",
// 					API:    "topic",
// 					Query:  map[string]interface{}{"name": "t_id"},
// 				},
// 			},
// 			want: jsonResponse{
// 				Err:    "",
// 				Status: http.StatusOK,
// 				Data:   nil,
// 			},
// 		},
// 		{
// 			name: "err - add topic - topic id empty",
// 			args: args{
// 				request: tcpTestRequest{
// 					Method: "POST",
// 					API:    "topic",
// 					Query:  map[string]interface{}{"name": ""},
// 				},
// 			},
// 			want: jsonResponse{
// 				Err:    "topicID cannot be empty",
// 				Status: http.StatusBadRequest,
// 				Data:   nil,
// 			},
// 		},
// 		{
// 			name: "err - add topic - duplicate id",
// 			args: args{
// 				request: tcpTestRequest{
// 					Method: "POST",
// 					API:    "topic",
// 					Query:  map[string]interface{}{"name": "t_id"},
// 				},
// 			},
// 			addtopic: true,
// 			want: jsonResponse{
// 				Err:    "topic ID already exists",
// 				Status: http.StatusConflict,
// 				Data:   nil,
// 			},
// 		},
// 		{
// 			name: "ok - get topics and subscribers",
// 			args: args{
// 				request: tcpTestRequest{
// 					Method: "GET",
// 					API:    "topic",
// 					Query:  nil,
// 				},
// 			},
// 			addtopic: true,
// 			want: jsonResponse{
// 				Err:    "",
// 				Status: http.StatusOK,
// 				Data: map[string]interface{}{
// 					"t_id": map[string]interface{}{},
// 				},
// 			},
// 		},
// 		{
// 			name: "err - get topics and subscribers - no stored topics",
// 			args: args{
// 				request: tcpTestRequest{
// 					Method: "GET",
// 					API:    "topic",
// 					Query:  nil,
// 				},
// 			},
// 			want: jsonResponse{
// 				Err:    "no topics stored",
// 				Status: http.StatusNotFound,
// 				Data:   nil,
// 			},
// 		},
// 		{
// 			name: "ok - send message to topic",
// 			args: args{
// 				request: tcpTestRequest{
// 					Method: "POST",
// 					API:    "message",
// 					Query:  map[string]interface{}{"topic": "t_id", "body": "m_body"},
// 				},
// 			},
// 			addtopic:       true,
// 			addSubscriber:  true,
// 			subscribeTopic: true,
// 			want: jsonResponse{
// 				Err:    "",
// 				Status: http.StatusOK,
// 				Data:   nil,
// 			},
// 		},
// 		{
// 			name: "err - send message to topic - topic id empty",
// 			args: args{
// 				request: tcpTestRequest{
// 					Method: "POST",
// 					API:    "message",
// 					Query:  map[string]interface{}{"topic": "", "body": "m_body"},
// 				},
// 			},
// 			addtopic:      true,
// 			addSubscriber: true,
// 			want: jsonResponse{
// 				Err:    "message topic empty",
// 				Status: http.StatusBadRequest,
// 				Data:   nil,
// 			},
// 		},
// 		{
// 			name: "err - send message to topic - topic body empty",
// 			args: args{
// 				request: tcpTestRequest{
// 					Method: "POST",
// 					API:    "message",
// 					Query:  map[string]interface{}{"topic": "t_id", "body": ""},
// 				},
// 			},
// 			addtopic:      true,
// 			addSubscriber: true,
// 			want: jsonResponse{
// 				Err:    "message body empty",
// 				Status: http.StatusBadRequest,
// 				Data:   nil,
// 			},
// 		},
// 		{
// 			name: "err - send message to topic - no topic subscribers",
// 			args: args{
// 				request: tcpTestRequest{
// 					Method: "POST",
// 					API:    "message",
// 					Query:  map[string]interface{}{"topic": "t_id", "body": "m_body"},
// 				},
// 			},
// 			addtopic: true,
// 			want: jsonResponse{
// 				Err:    "no subscribers to topic",
// 				Status: http.StatusNotFound,
// 				Data:   nil,
// 			},
// 		},
// 		// {
// 		// 	name: "ok - create new subscriber",
// 		// 	args: args{
// 		// 		request: tcpTestRequest{
// 		// 			Method: "POST",
// 		// 			API:    "broker",
// 		// 			Query:  nil,
// 		// 		},
// 		// 	},
// 		// 	mockSubscriberID: true,
// 		// 	want: jsonResponse{
// 		// 		Err:    "",
// 		// 		Status: http.StatusOK,
// 		// 		Data: map[string]interface{}{
// 		// 			"s_id": map[string]interface{}{
// 		// 				"ID":     "s_id",
// 		// 				"Topics": map[string]interface{}{},
// 		// 				"Active": true,
// 		// 			},
// 		// 		},
// 		// 	},
// 		// },
// 	}

// 	for _, tt := range tests {
// 		logger := logger.NewLogger()
// 		metrics := metrics.NewMetrics(logger)
// 		handlers := NewHandlers(logger)
// 		tcp := NewTCP(logger, metrics, handlers)

// 		logger.StartNoopLogger()
// 		tcp.Start()
// 		metrics.Start()

// 		t.Run(tt.name, func(t *testing.T) {
// 			conn, err := net.Dial("tcp", ":8181")
// 			if err != nil {
// 				t.Errorf("dial udp error: %v", err)
// 			}

// 			defer conn.Close()

// 			if tt.addtopic {
// 				tcp.handlers.br.AddTopic("t_id")
// 			}

// 			if tt.addSubscriber {
// 				sub := subscriber.NewSubscriber(logger)
// 				sub.ID = "s_id"

// 				tcp.handlers.br.AddSubscriber(sub)
// 				if tt.subscribeTopic {
// 					tcp.handlers.br.SubscribeToTopic(sub, "t_id")
// 				}

// 			}

// 			out, err := json.Marshal(tt.args.request)
// 			if err != nil {
// 				t.Errorf("marshal data error: %v", err)
// 			}

// 			_, err = conn.Write(out)
// 			if err != nil {
// 				t.Errorf("tcp write error: %v", err)
// 			}

// 			buf := make([]byte, 1024)
// 			n, err := conn.Read(buf)
// 			if err != nil {
// 				t.Errorf("tcp read error: %v", err)
// 			}

// 			var got jsonResponse
// 			err = json.Unmarshal(buf[0:n], &got)
// 			if err != nil {
// 				t.Errorf("json unmarshal error: %v", err)
// 			}

// 			// if tt.mockSubscriberID {
// 			// 	sub, ok := got.Data.(*subscriber.Subscriber)
// 			// 	if ok {
// 			// 		sub.ID = "s_id"
// 			// 	}

// 			// 	submap := map[string]interface{}

// 			// 	got.Data = map[string]interface{}{sub.ID: sub}
// 			// 	fmt.Println("got:::::: ", got)
// 			// }

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("TCPServer got = %v, want %v", got, tt.want)
// 			}

// 			conn.Close()
// 		})
// 		tcp.Stop()
// 	}
// }
