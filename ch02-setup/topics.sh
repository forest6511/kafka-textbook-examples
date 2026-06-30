#!/usr/bin/env bash
# ch02-setup: 第2章の CLI 疎通を再現するスクリプト。
#
# 事前にルートで `docker compose up -d` を実行しておくこと。
# 本文の kafka-topics / console producer / console consumer の操作を
# まとめて確認できる。
set -euo pipefail

BROKER=broker
KAFKA=/opt/kafka/bin
BOOTSTRAP=localhost:9092

run() { docker compose exec "$BROKER" "$@"; }

echo "== 1. トピック一覧(作成前は空) =="
run "$KAFKA/kafka-topics.sh" --bootstrap-server "$BOOTSTRAP" --list || true

echo "== 2. トピック orders を作成 =="
run "$KAFKA/kafka-topics.sh" --bootstrap-server "$BOOTSTRAP" \
    --create --topic orders --partitions 3 --replication-factor 1 || true

echo "== 3. orders の詳細 =="
run "$KAFKA/kafka-topics.sh" --bootstrap-server "$BOOTSTRAP" \
    --describe --topic orders

echo "== 4. メッセージを投入(hello / world / kafka) =="
printf 'hello\nworld\nkafka\n' | docker compose exec -T "$BROKER" \
    "$KAFKA/kafka-console-producer.sh" --bootstrap-server "$BOOTSTRAP" \
    --topic orders

echo "== 5. 最初から読み出す(5秒で終了) =="
docker compose exec -T "$BROKER" \
    "$KAFKA/kafka-console-consumer.sh" --bootstrap-server "$BOOTSTRAP" \
    --topic orders --from-beginning --timeout-ms 5000 || true

echo "== done =="
