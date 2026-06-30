#!/usr/bin/env bash
# コンシューマラグを発生させて監視で捉える。
# 5 件を消費して止め、さらに 4 件 produce すると LAG=4 になる。
set -euo pipefail

echo "=== 5 件を produce ==="
for i in 1 2 3 4 5; do echo "order-$i"; done | \
  docker compose exec -T broker-1 /opt/kafka/bin/kafka-console-producer.sh \
  --bootstrap-server broker-1:19092 --topic orders

echo "=== group app で消費して止める ==="
docker compose exec -T broker-1 /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server broker-1:19092 --topic orders --group app \
  --from-beginning --timeout-ms 5000 || true

echo "=== 消費を止めた状態でさらに 4 件 produce ==="
for i in 6 7 8 9; do echo "order-$i"; done | \
  docker compose exec -T broker-1 /opt/kafka/bin/kafka-console-producer.sh \
  --bootstrap-server broker-1:19092 --topic orders

echo "=== ラグを確認(LAG 列が増えている) ==="
docker compose exec -T broker-1 /opt/kafka/bin/kafka-consumer-groups.sh \
  --bootstrap-server broker-1:19092 --describe --group app
