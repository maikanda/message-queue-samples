package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func main() {
	consumer, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers":               "127.0.0.1:9092",
			"group.id":                        "myGroup",
			"go.application.rebalance.enable": true, // 再バランシングを有効化
		},
	)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	err = consumer.SubscribeTopics([]string{"test-topic"}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		consumer.Close()
		os.Exit(1)
	}

	run := true
	fmt.Println("Start consume")
	msg_count := 0
	for run {
		// Pollはconsumerが受け取ったメッセージ・イベントの双方を返す
		event := consumer.Poll(12000)
		switch e := event.(type) {
		// AssignedPartitions, RevokedPartitionsはリバランス時のイベント。consumerとpartitionの関連付けが更新されるタイミングで発生する
		case kafka.AssignedPartitions:
			fmt.Printf("AssignedPartitions: %v\n", e)
			err = consumer.Assign(e.Partitions)
			if err != nil {
				fmt.Printf("Assign error: %v\n", err)
			}
		case kafka.RevokedPartitions:
			fmt.Printf("RevokedPartitions %v\n", e)
			err = consumer.Unassign()
			if err != nil {
				fmt.Printf("UnAssign error: %v\n", err)
			}
		case *kafka.Message:
			msg_count += 1
			if msg_count%3 == 0 {
				go func() {
					offsets, err := consumer.Commit()
					if err != nil {
						fmt.Printf("Failed to commit error: %v offset: %v", err, offsets)
					}
				}()
			}
			fmt.Printf("%% Message on %s:%s\n", e.TopicPartition, string(e.Value))

		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
			fmt.Printf("Event %v\n", event)
		}
	}
	consumer.Close()
}
