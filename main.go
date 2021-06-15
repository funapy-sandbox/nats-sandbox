package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	subscription = "ORDERS.scratch"
	consumerName = "nats-consumer"
)

func run(ctx context.Context) error {
	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return fmt.Errorf("failed to apply: %w", err)
	}

	for i := 0; i < 10; i++ {
		js.PublishAsync(subscription, []byte("hello"))
	}
	select {
	case <-js.PublishAsyncComplete():
	case <-time.After(5 * time.Second):
		return errors.New("publish did not resolve in time")
	}
	// Simple Pull Consumer
	sub, err := js.PullSubscribe(subscription, consumerName)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	msgs, err := sub.Fetch(10)
	if err != nil {
		return fmt.Errorf("failed to fetch messages: %w", err)
	}

	for _, msg := range msgs {
		fmt.Println(string(msg.Data))
		if err := msg.Ack(); err != nil {
			log.Printf("failed to ack: %s\n", err.Error())
		}
	}

	sub.Unsubscribe()

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("error occurs: %s", err.Error())
	}
}
