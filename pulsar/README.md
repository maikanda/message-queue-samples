# 概要
PulsarのProducerとConsumerのサンプルです

# 実行方法

1. Pulsarの起動とテナント・ネームスペース・トピックの作成
```
docker compose up -d pulsar
docker compose exec pulsar bin/pulsar-admin tenants create my-tenant
docker compose exec pulsar bin/pulsar-admin namespaces create my-tenant/my-namespace
docker compose exec pulsar bin/pulsar-admin topics create-partitioned-topic persistent://my-tenant/my-namespace/my-topic -p 3
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

4. UIにログインするユーザーの作成
```
TOKEN=curl "http://localhost:7750/pulsar-manager/csrf-token"
curl -H "X-XSRF-TOKEN: $TOKEN" -H "Cookie: XSRF-TOKEN=$TOKEN" -H 'Content-Type: application/json' \
    -X PUT http://localhost:7750/pulsar-manager/users/superuser \
    -d '{"name": "admin", "password": "apachepulsar", "description": "test", "email": "username@test.org"}'
```

5. http://localhost:9527 を開き、name: admin pass: apachepulsar を入力しログイン