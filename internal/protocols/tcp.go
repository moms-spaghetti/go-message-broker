package protocols

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"task2/internal/logger"
	"task2/internal/message"
	"task2/internal/metrics"
	"task2/internal/subscriber"
)

var ErrApiForbidden = errors.New("api forbidden")
var ErrMethodForbidden = errors.New("method forbidden")
var ErrMessageRetryLimit = errors.New("message retries exceeded")

const (
	tcpaddr    = ":8181"
	tcpnetwork = "tcp"
)

type TCPServer struct {
	listener   net.Listener
	conns      map[string]net.Conn
	done       chan struct{}
	logger     *logger.Logger
	metrics    *metrics.Metrics
	handlers   *handlers
	connsMutex *sync.RWMutex
}

func NewTCP(
	logger *logger.Logger,
	metrics *metrics.Metrics,
	handlers *handlers,
) *TCPServer {

	lis, err := net.Listen(tcpnetwork, tcpaddr)
	if err != nil {
		panic(err)
	}

	return &TCPServer{
		listener:   lis,
		conns:      make(map[string]net.Conn),
		done:       make(chan struct{}),
		logger:     logger,
		metrics:    metrics,
		handlers:   handlers,
		connsMutex: &sync.RWMutex{},
	}
}

func (ts *TCPServer) Start() {
	go func() {
		for {
			log.Printf("tcp listning on %s", tcpaddr)
			conn, err := ts.listener.Accept()

			if err != nil {
				select {
				case <-ts.done:
					if len(ts.conns) > 0 {
						ts.logger.Log("closing active conns")

						for n, c := range ts.conns {
							ts.logger.Log(fmt.Sprintf("active conn %s closed", n))
							c.Close()
						}
					}
					return
				default:
					log.Printf("listener error: %v", err)

					return
				}
			}

			go ts.router(conn)
		}
	}()
}

func (ts *TCPServer) Stop() {
	close(ts.done)
	if err := ts.listener.Close(); err != nil {
		log.Printf("TCP listener close err: %v", err)
	}
	log.Print("TCP shutdown ok")
}

func (ts *TCPServer) router(conn net.Conn) {
	var (
		err error
		req jsonRequest
		res []byte
	)

	connID := createConnID()
	ts.addConn(conn, connID)

	err = json.NewDecoder(conn).Decode(&req)

	if err == nil {
		ts.logger.Log("api: " + req.API)
		ts.metrics.LogMetrics(req.API)
		switch req.API {
		case createSubscriber:
			res = ts.handlers.createSubscriber(req)
		case subscribeToTopic:
			res = ts.handlers.subscribeToTopic(req)
		case createTopic:
			res = ts.handlers.createTopic(req)
		case publishMessage:
			res = ts.handlers.publishMessage(req)
		case getMessages:
			res = ts.getMessages(req, conn)
		default:
			res = BuildJsonResponse(ErrApiForbidden, nil, ts.logger)
		}
	}

	conn.Write(res)
	conn.Close()
	ts.removeConn(connID)
}

func (ts *TCPServer) getMessages(req jsonRequest, conn net.Conn) []byte {
	var (
		err   error
		sub   *subscriber.Subscriber
		count int
		msg   message.Message
		out   []byte
		in    struct {
			Data struct {
				Message message.Message `json:"message"`
			}
		}
		r   int
		buf []byte
		n   int
	)

	var query reqMessagingQuery
	err = json.Unmarshal(req.Query, &query)

	if err == nil {
		sub, err = ts.handlers.br.GetSubscriber(query.SubID)
	}

	if err == nil {
		err = sub.ValidateTopicAndMessageCount(query.Topic)
	}

	if err == nil {
		count = sub.GetMessageCount(query.Topic)

		for {
			msg = sub.GetNextMessage(query.Topic)
			out = BuildJsonResponse(
				err,
				struct {
					Message message.Message `json:"message"`
					Count   int             `json:"count"`
				}{
					Message: msg,
					Count:   count,
				},
				ts.logger,
			)
			_, err = conn.Write(out)
			buf = make([]byte, 1024)
			if err == nil {
				n, err = conn.Read(buf)
			}
			if err == nil {
				err = json.Unmarshal(buf[:n], &in)
			}
			if err == nil {
				if !in.Data.Message.Done {
					r++
				} else {
					err = sub.DeleteMessage(query.Topic, in.Data.Message.ID, ts.logger)
					count = sub.GetMessageCount(query.Topic)
					r = 0
				}
			}
			if r == 3 {
				err = errors.New("retries exceeded for message id " + in.Data.Message.ID)
				break
			}
			if err != nil || count == 0 {
				break
			}
		}
	}

	return BuildJsonResponse(err, nil, ts.logger)
}

// conn clean up
func (ts *TCPServer) addConn(conn net.Conn, connID string) {
	ts.logger.Log("conn added: " + connID)

	ts.connsMutex.Lock()
	defer ts.connsMutex.Unlock()

	ts.conns[connID] = conn
}

func (ts *TCPServer) removeConn(connID string) {
	ts.logger.Log("conn removed: " + connID)

	ts.connsMutex.Lock()
	defer ts.connsMutex.Unlock()

	delete(ts.conns, connID)
}

// horrible dirty ID creation
func createConnID() string {
	id := make([]byte, 4)
	rand.Read(id)

	return fmt.Sprintf("%X", id[0:4])
}
