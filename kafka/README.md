# 概要
KafkaのProducerとConsumerのサンプルです

# 実行方法

1. Kafkaの起動とトピックの作成
```
docker compose up -d kafka
docker compose exec kafka ./kafka-topics.sh --create --topic test-topic --bootstrap-server localhost:9092 --partitions 3
```

2. Producerの起動
```
docker compose up -d producer
docker compose exec producer ./exec.sh
```   

3. Consumerの起動
```
docker compose up -d consumer
docker compose exec consumer go run main.go
```