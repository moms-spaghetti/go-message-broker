package protocols

import (
	"encoding/json"
	"log"
	"net"
	"task2/internal/logger"
	"task2/internal/message"
	"task2/internal/metrics"
	"task2/internal/subscriber"
)

const (
	udpnetwork = "udp"
	udpport    = 9001
	zone       = ""
	bufsize    = 1024
)

type UDPServer struct {
	laddr    *net.UDPAddr
	done     chan struct{}
	logger   *logger.Logger
	metrics  *metrics.Metrics
	handlers *handlers
}

func NewUDP(
	logger *logger.Logger,
	metrics *metrics.Metrics,
	handlers *handlers,
) *UDPServer {
	laddr := &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: udpport,
		Zone: zone,
	}

	return &UDPServer{
		laddr:    laddr,
		done:     make(chan struct{}),
		logger:   logger,
		metrics:  metrics,
		handlers: handlers,
	}
}

func (us *UDPServer) Start() {
	go func() {
		conn, err := net.ListenUDP(udpnetwork, us.laddr)
		if err != nil {
			panic(err)
		}

		log.Printf("udp listening on %s", us.laddr)

		for {
			select {
			case <-us.done:
				if err := conn.Close(); err != nil {
					log.Printf("UDP listener close err: %v", err)
				}

				log.Print("UDP shutdown ok")
				return
			default:
				buf := make([]byte, bufsize)
				n, raddr, err := conn.ReadFromUDP(buf)
				if err != nil {
					log.Printf("startUDP error: %v", err)
				}

				go us.router(conn, buf[:n], raddr)
			}
		}
	}()
}

func (us *UDPServer) Stop() {
	close(us.done)

}

func (us *UDPServer) router(conn *net.UDPConn, buf []byte, raddr *net.UDPAddr) {
	var (
		err error
		req jsonRequest
		res []byte
	)

	err = json.Unmarshal(buf, &req)

	if err == nil {
		us.logger.Log("api: " + req.API)
		us.metrics.LogMetrics(req.API)
		switch req.API {
		case createSubscriber:
			res = us.handlers.createSubscriber(req)
		case subscribeToTopic:
			res = us.handlers.subscribeToTopic(req)
		case createTopic:
			res = us.handlers.createTopic(req)
		case publishMessage:
			res = us.handlers.publishMessage(req)
		case udpGetMessageCount, udpGetNextMessage, udpCompleteMessage:
			res = us.getMessages(req, conn)
		default:
			res = BuildJsonResponse(ErrApiForbidden, nil, us.logger)
		}
	}

	conn.WriteTo(res, raddr)
}

func (us *UDPServer) getMessages(req jsonRequest, conn net.Conn) []byte {
	var (
		query reqMessagingQuery
		err   error
		sub   *subscriber.Subscriber
		res   []byte
	)
	err = json.Unmarshal(req.Query, &query)

	switch req.API {
	case udpGetMessageCount:
		var data struct {
			Count int `json:"count"`
		}

		if err == nil {
			sub, err = us.handlers.br.GetSubscriber(query.SubID)
		}

		if err == nil {
			err = sub.ValidateTopicAndMessageCount(query.Topic)
		}

		if err == nil {
			data.Count = sub.GetMessageCount(query.Topic)
			res = BuildJsonResponse(
				err,
				data,
				us.logger,
			)
		}
	case udpGetNextMessage:
		var data struct {
			Message message.Message `json:"message"`
		}
		if err == nil {
			sub, err = us.handlers.br.GetSubscriber(query.SubID)
		}

		if err == nil {
			err = sub.ValidateTopicAndMessageCount(query.Topic)
		}

		if err == nil {
			data.Message = sub.GetNextMessage(query.Topic)
			res = BuildJsonResponse(
				err,
				data,
				us.logger,
			)
		}
	case udpCompleteMessage:
		var message message.Message
		if err == nil {
			err = json.Unmarshal(req.Payload, &message)
		}

		if err == nil {
			sub, err = us.handlers.br.GetSubscriber(query.SubID)
		}

		if err == nil {
			err = sub.ValidateTopicAndMessageCount(query.Topic)
		}

		if err == nil {
			if message.Done {
				err = sub.DeleteMessage(query.Topic, message.ID, us.logger)
			}
			res = BuildJsonResponse(
				err,
				nil,
				us.logger,
			)
		}
	}

	return res
}
