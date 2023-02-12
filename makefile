.PHONY: run runc

run:
		go run cmd/broker/main.go

runc:
		go run cmd/client/main.go