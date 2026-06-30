#!/usr/bin/env bash
# ch07-transactions: consume-transform-produce のトランザクション処理を end-to-end で確認。
#
# orders を read_committed で読み、value に "-processed" を付けて orders-processed へ
# トランザクションで原子的に書き込む。読んだオフセットのコミットと produce が
# 1 つのトランザクションでまとまる(exactly-once)。
#
# 事前にルートで `docker compose up -d` を実行しておくこと。
set -euo pipefail

BROKER=broker
KAFKA=/opt/kafka/bin
BOOTSTRAP=localhost:9092

echo "== orders / orders-processed を作り直す =="
for T in orders orders-processed; do
  docker compose exec -T "$BROKER" "$KAFKA/kafka-topics.sh" \
      --bootstrap-server "$BOOTSTRAP" --delete --topic "$T" 2>/dev/null || true
done
sleep 2
for T in orders orders-processed; do
  docker compose exec -T "$BROKER" "$KAFKA/kafka-topics.sh" \
      --bootstrap-server "$BOOTSTRAP" \
      --create --topic "$T" --partitions 3 --replication-factor 1
done

echo "== ch04 Producer で orders に 3 件送る =="
go run ./ch04-producer

echo "== ch07 トランザクション処理(orders -> orders-processed) =="
go run ./ch07-transactions

echo "== orders-eos グループのラグ(LAG=0 なら全件処理済み) =="
docker compose exec -T "$BROKER" "$KAFKA/kafka-consumer-groups.sh" \
    --bootstrap-server "$BOOTSTRAP" --describe --group orders-eos

echo "== orders-processed を read_committed で読む =="
docker compose exec -T "$BROKER" "$KAFKA/kafka-console-consumer.sh" \
    --bootstrap-server "$BOOTSTRAP" --topic orders-processed \
    --from-beginning --isolation-level read_committed --timeout-ms 6000 \
    --property print.key=true --property print.partition=true || true

echo "== done =="
