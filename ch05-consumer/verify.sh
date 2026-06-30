#!/usr/bin/env bash
# ch05-consumer: ch04 Producer で送り、ch05 Consumer で受け取る end-to-end 確認。
#
# 事前にルートで `docker compose up -d` を実行しておくこと。
set -euo pipefail

BROKER=broker
KAFKA=/opt/kafka/bin
BOOTSTRAP=localhost:9092
TOPIC=orders
GROUP=orders-processor

echo "== orders を作り直してクリーンな状態にする =="
docker compose exec -T "$BROKER" "$KAFKA/kafka-topics.sh" \
    --bootstrap-server "$BOOTSTRAP" --delete --topic "$TOPIC" 2>/dev/null || true
sleep 2
docker compose exec -T "$BROKER" "$KAFKA/kafka-topics.sh" \
    --bootstrap-server "$BOOTSTRAP" \
    --create --topic "$TOPIC" --partitions 3 --replication-factor 1

echo "== ch04 Producer でメッセージを送る =="
go run ./ch04-producer

echo "== ch05 Consumer で受け取る（7 秒で停止） =="
go build -o /tmp/ch05c ./ch05-consumer
/tmp/ch05c &
PID=$!
sleep 7
kill -INT "$PID" 2>/dev/null || true
wait "$PID" 2>/dev/null || true

echo "== グループのオフセット/ラグを確認（LAG=0 なら全件処理済み） =="
docker compose exec -T "$BROKER" "$KAFKA/kafka-consumer-groups.sh" \
    --bootstrap-server "$BOOTSTRAP" --describe --group "$GROUP"

echo "== done =="
