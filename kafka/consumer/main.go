package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	//ctx, cancel := context.WithCancel(context.Background())

	// ConfigMapの詳細はこちら https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md
	consumer1, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers":               "kafka:9092",
			"group.id":                        "myGroup",
			"go.application.rebalance.enable": true, // 再バランシングを有効化
		},
	)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	//consumer2, err := kafka.NewConsumer(
	//	&kafka.ConfigMap{
	//		"bootstrap.servers":               "kafka:9092",
	//		"group.id":                        "myGroup",
	//		"go.application.rebalance.enable": true, // 再バランシングを有効化
	//	},
	//)
	//if err != nil {
	//	fmt.Printf("Failed to create consumer: %s\n", err)
	//	os.Exit(1)
	//}

	err = consumer1.SubscribeTopics([]string{"test-topic"}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		consumer1.Close()
		os.Exit(1)
	}

	//err = consumer2.SubscribeTopics([]string{"test-topic"}, nil)
	//if err != nil {
	//	fmt.Printf("Failed to subscribe to topic: %s\n", err)
	//	consumer2.Close()
	//	os.Exit(1)
	//}

	defer func() {
		if err = consumer1.Close(); err != nil {
			fmt.Printf("Failed to close consumer1: %s\n", err)
		}
		//if err = consumer2.Close(); err != nil {
		//	fmt.Printf("Failed to close consumer2: %s\n", err)
		//}
	}()

	for {
		m, _ := consumer1.ReadMessage(10)
		fmt.Printf("Message %v\n", m)
	}
	//go Consume(consumer1, ctx)
	//go Consume(consumer2, ctx)

	<-stopChan
	//cancel()
}

func Consume(consumer *kafka.Consumer, ctx context.Context) {
	run := true
	fmt.Println("Start consume")
	for run {
		select {
		case <-ctx.Done():
			run = false
			return
		default:
			// Pollはconsumerが受け取ったメッセージ・イベントの双方を返す
			event := consumer.Poll(0)
			switch e := event.(type) {
			// AssignedPartitions, RevokedPartitionsはリバランス時のイベント。consumerとpartitionの関連付けが更新されるタイミングで発生する
			case kafka.AssignedPartitions:
				fmt.Printf("AssignedPartitions: %v\n", e)
				err := consumer.Assign(e.Partitions)
				if err != nil {
					fmt.Printf("Assign error: %v\n", err)
				}
			case kafka.RevokedPartitions:
				fmt.Printf("RevokedPartitions %v\n", e)
				err := consumer.Unassign()
				if err != nil {
					fmt.Printf("UnAssign error: %v\n", err)
				}
			case *kafka.Message:
				// Consumerのenable.auto.commitによって自動的にcommitされる
				// デフォルトで5000msごとにcommitされる
				fmt.Printf("%% Consumer %v Message on %s:%s\n", consumer.String(), e.TopicPartition, string(e.Value))

			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			default:
			}
		}
	}
}
