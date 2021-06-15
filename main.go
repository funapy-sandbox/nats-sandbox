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
	streamName    = "sandbox-stream"
	streamSubject = "sandbox-stream.*"
	subject       = "sandbox-stream.test-subect"
	consumerName  = "NATS-CONSUMER"
)

func run(ctx context.Context) error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return fmt.Errorf("failed to apply: %w", err)
	}
	info, err := js.StreamInfo(streamName)
	if err != nil {
		log.Printf("failed to get stream info: %s\n", err)
	}

	if info == nil {
		log.Printf("create stream")
		_, err := js.AddStream(&nats.StreamConfig{
			Name: streamName,
			Subjects: []string{
				streamSubject,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to add stream :%w", err)
		}
	}

	_, err = js.Publish(subject, []byte("aaa"))
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	for i := 0; i < 10; i++ {
		js.PublishAsync(subject, []byte("hello"))
	}
	select {
	case <-js.PublishAsyncComplete():
		log.Println("publish async complete")
	case <-time.After(5 * time.Second):
		return errors.New("publish did not resolve in time")
	}

	sub, err := js.PullSubscribe(subject, consumerName)
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
