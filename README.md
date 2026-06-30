# kafka-textbook-examples

「**Apache Kafka の教科書 — Go で学ぶ Apache Kafka ― Producer・Consumer から KRaft 運用まで**」（森川 陽介 著）のサンプルコード集です。

本書の各章で使う Go コード（franz-go）・`docker compose` 構成・Kafka CLI コマンドを、そのまま動かせる形で収録しています。すべて **Apache Kafka 4.x（KRaft モード）** と **franz-go v1.21+** で動作を確認しています。本書のターミナル出力は、ここの構成を実機で起動して取得した値です。

## 前提

- Docker / Docker Compose（Kafka を KRaft モードで起動）
- Go 1.26 以降（franz-go クライアント）

## 共通の Kafka を起動する

リポジトリのルートに、KRaft 単一ノードの `docker-compose.yml` があります。ほとんどの章はこれを使います。

```bash
docker compose up -d        # localhost:9092 で起動
docker compose ps           # 状態確認
docker compose down -v      # 停止 + データ削除
```

CLI はコンテナ内のものを使えます。

```bash
docker compose exec broker /opt/kafka/bin/kafka-topics.sh \
    --bootstrap-server localhost:9092 --list
```

## 構成（各章ディレクトリ）

コードが登場する章ごとに `chNN-<topic>/` ディレクトリがあり、その章のコードと実行手順（`README.md`）が入っています。

- ch02-setup — Docker で KRaft クラスタを起動し、franz-go から最小疎通する
- ch04-producer — franz-go で Producer を書く
- ch05-consumer — franz-go で Consumer を書く
- ch07-transactions — 配信保証とトランザクション（consume-transform-produce の exactly-once）
- ch08-schema — Schema Registry 連携
- ch09-operations — 運用（複数ブローカー構成・監視）
- ch10-client-selection — franz-go / confluent-kafka-go / segmentio の書き比べ
- ch11-systematize — テスト（kfake）・CI・再現性

> 第1章（Kafka をいつ使うかの判断章）・第3章（中核概念）・第6章（イベント駆動の設計章）は実行するコードが中心ではないため、ディレクトリはありません（第3章で使うコマンドは ch03-concepts に、第6章の Producer/Consumer は ch04-producer / ch05-consumer を役割で読み替えます）。

## バージョン

- Apache Kafka 4.x（`apache/kafka` 公式イメージ・KRaft モード、ZooKeeper 不要）
- franz-go v1.21+（Pure Go・cgo 不要）
