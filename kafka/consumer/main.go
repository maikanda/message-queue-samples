package main

import (
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

fun main() {
	consumer, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers":    "host1:9092,host2:9092",
			"group.id":             "foo",
			"auto.offset.reset":    "smallest"
		}
	)
}

