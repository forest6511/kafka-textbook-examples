#!/usr/bin/env bash
# ch10-clients/verify.sh — 3クライアントで同一トピックに疎通する
#
# 前提: ルートで `docker compose up -d` 済み。
#   bash ch10-clients/verify.sh
set -euo pipefail

cd "$(dirname "$0")/.."

echo "== トピック作成 =="
docker compose exec broker /opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --if-not-exists --topic ch10-compare \
  --partitions 1 --replication-factor 1

echo "== franz-go (Pure Go, CGO_ENABLED=0) =="
CGO_ENABLED=0 go run ./ch10-clients/franz

echo "== segmentio (Pure Go, CGO_ENABLED=0) =="
CGO_ENABLED=0 go run ./ch10-clients/segmentio

echo "== confluent (cgo 必須: CGO_ENABLED=0 ではビルドできない) =="
if CGO_ENABLED=0 go build -tags confluent -o /dev/null ./ch10-clients/confluent 2>/dev/null; then
  echo "想定外: CGO_ENABLED=0 でビルドできてしまった"
else
  echo "想定どおり CGO_ENABLED=0 ではビルド失敗(librdkafka 未リンク)"
fi
echo "-- cgo を有効にして実行 --"
CGO_ENABLED=1 go run -tags confluent ./ch10-clients/confluent

echo "== トピックの全メッセージ =="
docker compose exec broker /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 --topic ch10-compare \
  --from-beginning --timeout-ms 4000 2>/dev/null | sort
