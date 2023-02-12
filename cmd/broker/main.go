package main

import (
	"os"
	"os/signal"
	"syscall"
	"task2/internal/logger"
	"task2/internal/metrics"
	"task2/internal/protocols"
)

func main() {
	logger := logger.NewLogger()

	metrics := metrics.NewMetrics(logger)
	handlers := protocols.NewHandlers(logger)

	tcp := protocols.NewTCP(logger, metrics, handlers)
	udp := protocols.NewUDP(logger, metrics, handlers)

	starts := []func(){
		logger.Start,
		tcp.Start,
		udp.Start,
		metrics.Start,
	}

	stops := []func(){
		tcp.Stop,
		udp.Stop,
		metrics.Stop,
		logger.Stop,
	}

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT)

	run(starts)

	<-wait
	defer close(wait)

	run(stops)
}

func run(fn []func()) {
	for _, f := range fn {
		f()
	}
}
