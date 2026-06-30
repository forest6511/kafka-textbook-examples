#!/usr/bin/env bash
# ch04-producer: franz-go の Producer で orders に送り、CLI で届いたか確認する。
#
# 事前にルートで `docker compose up -d` を実行しておくこと。
set -euo pipefail

BROKER=broker
KAFKA=/opt/kafka/bin
BOOTSTRAP=localhost:9092
TOPIC=orders

echo "== トピック $TOPIC を用意(3 パーティション、なければ作成) =="
docker compose exec -T "$BROKER" "$KAFKA/kafka-topics.sh" \
    --bootstrap-server "$BOOTSTRAP" \
    --create --if-not-exists --topic "$TOPIC" \
    --partitions 3 --replication-factor 1

echo "== franz-go の Producer を実行 =="
go run ./ch04-producer

echo "== console-consumer で届いたメッセージを確認 =="
docker compose exec -T "$BROKER" "$KAFKA/kafka-console-consumer.sh" \
    --bootstrap-server "$BOOTSTRAP" --topic "$TOPIC" \
    --from-beginning --timeout-ms 4000 \
    --formatter-property print.key=true 2>/dev/null || true

echo "== done =="
