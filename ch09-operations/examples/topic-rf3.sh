#!/usr/bin/env bash
# replication.factor=3 + min.insync.replicas=2 のトピックを作り、
# Leader/Replicas/Isr を確認する。3 ブローカー構成で実行する。
set -euo pipefail

docker compose exec -T broker-1 /opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server broker-1:19092 --create --topic orders \
  --partitions 3 --replication-factor 3 \
  --config min.insync.replicas=2

docker compose exec -T broker-1 /opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server broker-1:19092 --describe --topic orders
