package metrics

import (
	"log"
	"task2/internal/logger"

	"fmt"
)

const (
	// custom api metrics
	createSubscriber   = "createSubscriber"
	subscribeToTopic   = "subscribeToTopic"
	createTopic        = "createTopic"
	publishMessage     = "publishMessage"
	getMessages        = "getMessages"
	udpGetMessageCount = "udpGetMessageCount"
	udpGetNextMessage  = "udpGetNextMessage"
	udpCompleteMessage = "udpCompleteMessage"
)

type Metrics struct {
	metrics chan string
	done    chan struct{}
	stats   *stats
	logger  *logger.Logger
}

type stats struct {
	createSubscriber   int
	subscribeToTopic   int
	createTopic        int
	publishMessage     int
	getMessages        int
	udpGetMessageCount int
	udpGetNextMessage  int
	udpCompleteMessage int
	unknown            int
}

func NewMetrics(logger *logger.Logger) *Metrics {
	metrics := &Metrics{
		metrics: make(chan string),
		done:    make(chan struct{}),
		stats: &stats{
			createSubscriber:   0,
			subscribeToTopic:   0,
			createTopic:        0,
			publishMessage:     0,
			getMessages:        0,
			udpGetMessageCount: 0,
			udpGetNextMessage:  0,
			udpCompleteMessage: 0,
			unknown:            0,
		},
		logger: logger,
	}

	return metrics
}

func (m *Metrics) Start() {
	log.Print("metrics started")
	go func() {
		for {
			select {
			case stat := <-m.metrics:
				switch stat {
				case createSubscriber:
					m.stats.createSubscriber++
				case subscribeToTopic:
					m.stats.subscribeToTopic++
				case createTopic:
					m.stats.createTopic++
				case publishMessage:
					m.stats.publishMessage++
				case getMessages:
					m.stats.getMessages++
				case udpGetMessageCount:
					m.stats.udpGetMessageCount++
				case udpGetNextMessage:
					m.stats.udpGetNextMessage++
				case udpCompleteMessage:
					m.stats.udpCompleteMessage++
				default:
					m.stats.unknown++
				}

				m.PrintMetrics()
			case <-m.done:
				return
			}
		}
	}()
}

func (m *Metrics) LogMetrics(method string) {
	m.metrics <- method
}

func (m *Metrics) PrintMetrics() {
	out := fmt.Sprintf(`
createSubscriber: %d
subscribeToTopic: %d
createTopic: %d
publishMessage: %d
getMessages: %d
udpGetMessageCount: %d
udpGetNextMessage: %d
udpCompleteMessage: %d
unknown: %d`,
		m.stats.createSubscriber,
		m.stats.subscribeToTopic,
		m.stats.createTopic,
		m.stats.publishMessage,
		m.stats.getMessages,
		m.stats.udpGetMessageCount,
		m.stats.udpGetNextMessage,
		m.stats.udpCompleteMessage,
		m.stats.unknown)
	m.logger.Log(out)
}

func (m *Metrics) Stop() {
	close(m.metrics)
	close(m.done)
	log.Print("metrics shutdown ok")
}
