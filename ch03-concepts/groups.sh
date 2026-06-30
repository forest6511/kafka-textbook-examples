#!/usr/bin/env bash
# ch03-concepts: コンシューマグループのオフセット/ラグを観察する。
#
# 事前にルートで `docker compose up -d` を実行しておくこと。
# orders に 12 件のキー付きメッセージを投入し、ch03-demo グループで
# 全件読んだあと、グループの読み進み具合を表示する。
set -euo pipefail

BROKER=broker
KAFKA=/opt/kafka/bin
BOOTSTRAP=localhost:9092
TOPIC=orders
GROUP=ch03-demo

run() { docker compose exec "$BROKER" "$@"; }

echo "== トピック $TOPIC を用意(なければ作成) =="
run "$KAFKA/kafka-topics.sh" --bootstrap-server "$BOOTSTRAP" \
    --create --if-not-exists --topic "$TOPIC" \
    --partitions 3 --replication-factor 1

echo "== 12 件のキー付きメッセージを投入 =="
printf 'k1:a\nk2:b\nk1:c\nk3:d\nk2:e\nk1:f\nk3:g\nk2:h\nk1:i\nk3:j\nk2:k\nk1:l\n' \
  | docker compose exec -T "$BROKER" \
      "$KAFKA/kafka-console-producer.sh" --bootstrap-server "$BOOTSTRAP" \
      --topic "$TOPIC" \
      --property "parse.key=true" --property "key.separator=:"

echo "== $GROUP グループで全件読む(5秒で終了) =="
docker compose exec -T "$BROKER" \
    "$KAFKA/kafka-console-consumer.sh" --bootstrap-server "$BOOTSTRAP" \
    --topic "$TOPIC" --from-beginning --group "$GROUP" --timeout-ms 5000 \
    > /dev/null 2>&1 || true

echo "== グループのオフセット/ラグを表示 =="
run "$KAFKA/kafka-consumer-groups.sh" --bootstrap-server "$BOOTSTRAP" \
    --describe --group "$GROUP"

echo "== done =="
