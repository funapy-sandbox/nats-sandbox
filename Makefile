
NATS_CONTAINER_NAME := "nats-sandbox"
NATS_VERSION        := "2.2.6"

.PHONY: nats/server/start
nats/server/start:
	docker run --name $(NATS_CONTAINER_NAME) -p 4222:4222 -p 6222:6222 -p 8222:8222 nats:$(NATS_VERSION) -js

.PHONY: nats/server/stop
nats/server/stop:
	docker rm -f $(NATS_CONTAINER_NAME)

.PHONY: run
run:
	go run ./main.go
