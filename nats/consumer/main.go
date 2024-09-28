package main

import (
	"context"
	"fmt"
	nats "github.com/nats-io/nats.go"
	jetstream "github.com/nats-io/nats.go/jetstream"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		fmt.Printf("Failed to connect nats %v\n", err)
		os.Exit(1)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		fmt.Printf("Failed to establish stream %v\n", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// streamの作成 (または更新)
	stream, _ := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "stream-1",
		Subjects: []string{"sub.my-topic"},
	})

	// メッセージを受信するためには、streamにconsumerを作成する
	consumer1, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:      "consumer-1",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		fmt.Printf("Failed to add consumer %v\n", err)
		os.Exit(1)
	}

	consumer2, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:      "consumer-2",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		fmt.Printf("Failed to add consumer %v\n", err)
		os.Exit(1)
	}

	// メモ 以下の実装について
	// 1つのstreamに複数consumerを作成し、メッセージを並列でバッチ受信すると、重複して受信してしまう
	go Consume(consumer1)
	go Consume(consumer2)

	<-stopChan
	cancel()
}

func Consume(consumer jetstream.Consumer) {
	for {
		// 最大10件受信する
		messages, err := consumer.Fetch(10)
		if err != nil {
			fmt.Printf("Failed to get message error: %v\\n", err)
			continue
		}

		for message := range messages.Messages() {
			message.Ack()
			meta, _ := message.Metadata()
			fmt.Printf("Received a JetStream message via fetch: %s stream: %v consumer: %v domain: %v sequence %v\n", string(message.Data()), meta.Stream, meta.Consumer, meta.Domain, meta.Sequence)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
