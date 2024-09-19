package main

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"
	"time"
)

func main() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://pulsar:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}

	defer client.Close()
	// KeySharedを用いると、同一サブスクリプション名の場合、特定のキーのみ受け取るようになる
	consumer1, err := client.Subscribe(pulsar.ConsumerOptions{
		Name:             "Consumer-1",
		SubscriptionName: "sub",
		Topic:            "persistent://my-tenant/my-namespace/my-topic",
		Type:             pulsar.KeyShared,
	})
	consumer2, err := client.Subscribe(pulsar.ConsumerOptions{
		Name:             "Consumer-2",
		SubscriptionName: "sub",
		Topic:            "persistent://my-tenant/my-namespace/my-topic",
		Type:             pulsar.KeyShared,
	})
	if err != nil {
		log.Fatal(err)
	}

	run := true
	for run {
		msg1, err := consumer1.Receive(context.Background())
		if err != nil {
			fmt.Printf("Failed to receive message consumer1 %v\n", err)
		}
		fmt.Printf("Receive consumer1 key: %v value: %v\n", msg1.Key(), string(msg1.Payload()))
		msg2, err := consumer2.Receive(context.Background())
		if err != nil {
			fmt.Printf("Failed to receive message consumer1 %v\n", err)
		}
		fmt.Printf("Receive consumer2 key: %v value: %v\n", msg2.Key(), string(msg2.Payload()))
	}
}
