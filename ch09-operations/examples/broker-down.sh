#!/usr/bin/env bash
# broker-1 を止めて、リーダーが ISR 内の別ブローカーに移るのを確認する。
# その後 broker-1 を戻し、ISR が回復するのを見る(ローリング再起動の考え方)。
set -euo pipefail

echo "=== broker-1 を停止 ==="
docker compose stop broker-1
sleep 8

echo "=== describe(broker-2 経由)。Leader が移り、Isr から 1 が消える ==="
docker compose exec -T broker-2 /opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server broker-2:19092 --describe --topic orders

echo "=== broker-1 を再起動 ==="
docker compose start broker-1
sleep 10

echo "=== ISR が 1,2,3 に回復するのを確認 ==="
docker compose exec -T broker-2 /opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server broker-2:19092 --describe --topic orders
