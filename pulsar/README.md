# 概要
PulsarのProducerとConsumerのサンプルです

# 実行方法

1. Pulsarの起動とテナント・ネームスペース・トピックの作成
```
docker compose up -d pulsar
docker compose exec pulsar bin/pulsar-admin tenants create my-tenant
docker compose exec pulsar bin/pulsar-admin namespaces create my-tenant/my-namespace
docker compose exec pulsar bin/pulsar-admin topics create-partitioned-topic persistent://my-tenant/my-namespace/my-topic -p 5
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