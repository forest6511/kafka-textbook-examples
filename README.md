# kafka-textbook-examples

「**Apache Kafka の教科書 — Go で学ぶ Apache Kafka ― Producer・Consumer から KRaft 運用まで**」（森川 陽介 著）のサンプルコード集です。

リポジトリ: <https://github.com/forest6511/kafka-textbook-examples>

本書の各章で使う Go コード（franz-go）・`docker compose` 構成・Kafka CLI コマンドを、そのまま動かせる形で収録しています。すべて **Apache Kafka 4.3.1（KRaft モード）** と **franz-go v1.21.4** で動作を確認しています。本書のターミナル出力は、ここの構成を実機で起動して取得した値です。

## クイックスタート

```bash
git clone https://github.com/forest6511/kafka-textbook-examples.git
cd kafka-textbook-examples

docker compose up -d          # localhost:9092 で Kafka を起動
go run ./ch02-setup           # franz-go から最小疎通（"Kafka に接続できました" が出れば成功）
```

## 前提

- Docker / Docker Compose（Kafka を KRaft モードで起動）
- Go 1.26 以降（`go.mod` の要求バージョン。franz-go クライアントのビルドに使います）

## 共通の Kafka を起動する

リポジトリのルートに、KRaft 単一ノードの `docker-compose.yml` があります。**ほとんどの章はこれを使います**（例外は下記「章とディレクトリの対応」を参照）。

```bash
docker compose up -d        # localhost:9092 で起動
docker compose ps           # 状態確認（STATUS が Up になれば起動完了）
docker compose down -v      # 停止 + データ削除
```

CLI はコンテナ内のものを使えます。公式イメージ `apache/kafka` では CLI は `/opt/kafka/bin/` 配下にあり、**本文中で `kafka-topics` と書いているコマンドは、このパスの `kafka-topics.sh` を指します**。

```bash
docker compose exec broker /opt/kafka/bin/kafka-topics.sh \
    --bootstrap-server localhost:9092 --list
```

## 章とディレクトリの対応

コードが登場する章ごとに `chNN-<topic>/` ディレクトリがあり、その章のコードと実行手順（`README.md`）が入っています。

- [ch02-setup](ch02-setup/) — Docker で KRaft クラスタを起動し、franz-go から最小疎通する
- [ch03-concepts](ch03-concepts/) — 中核概念（パーティション・オフセット・コンシューマグループ）を CLI で観察する
- [ch04-producer](ch04-producer/) — franz-go で Producer を書く
- [ch05-consumer](ch05-consumer/) — franz-go で Consumer を書く
- [ch07-transactions](ch07-transactions/) — 配信保証とトランザクション（consume-transform-produce の exactly-once）
- [ch08-schema](ch08-schema/) — Schema Registry 連携（**この章専用の `docker-compose.yml` を使います**）
- [ch09-operations](ch09-operations/) — 運用（3 ブローカー構成・監視、**この章専用の `docker-compose.yml` を使います**）
- [ch10-clients](ch10-clients/) — franz-go / confluent-kafka-go / segmentio の書き比べ
- [ch11-systematize](ch11-systematize/) — テスト（kfake）・CI・再現性

ディレクトリがないのは **第1章**（Kafka をいつ使うかの判断章）と **第6章**（イベント駆動の設計章）です。どちらも実行するコードが中心ではないため、コードはありません（第6章の Producer / Consumer は ch04-producer / ch05-consumer を役割で読み替えます）。**第3章**は概念章ですが、CLI で観察するコマンドをまとめた `ch03-concepts/` を用意しています。

> **ch08 / ch09 を動かすときの注意**: この 2 つの章は、章ディレクトリ内の専用 `docker-compose.yml`（Schema Registry / 3 ブローカー）を使います。ルートの Kafka（9092）を起動したままだとポートが衝突するので、先に `docker compose down -v` でルートの Kafka を止めてから、各章の README に従って起動してください。

## うまく動かないとき

- **`go run` が `connection refused` になる** — ブローカーは起動直後、受け付けまで数秒かかります。`docker compose ps` で `STATUS` が `Up` か確認し、必要なら `docker compose logs broker` を見てください。
- **`bind: address already in use` / ポート 9092 が使えない** — 別の Kafka やプロセスが 9092 を専有しています。`docker compose down -v` で掃除するか、既存のプロセスを止めてください。
- **章を切り替えたら動かない** — 前の章の Kafka が残っている可能性があります。章を移るときは前の章で `docker compose down -v` を実行してから次に進んでください（特に ch08 / ch09 の専用構成とルート構成は同時に起動できません）。
- **confluent-kafka-go が `undefined: kafka.ConfigMap` でビルド失敗する** — このクライアントは librdkafka（C）に依存するため cgo が必須です。`CGO_ENABLED=1` でビルドしてください。詳しくは [ch10-clients/README.md](ch10-clients/) を参照。

## バージョン

- Apache Kafka 4.3.1（`apache/kafka` 公式イメージ・KRaft モード、ZooKeeper 不要）
- franz-go v1.21.4（Pure Go・cgo 不要）
- confluent-kafka-go v2.15.0（cgo 必須） / segmentio/kafka-go v0.4.51（比較用・ch10）
- Confluent Schema Registry 7.9.0（ch08）
- Go 1.26

（検証日: 2026-07-01）

## 書籍について

本書は Amazon Kindle / ペーパーバックで販売予定です（出版後、ここに商品ページのリンクを追加します）。

## ライセンス

MIT License. 本書のサンプルコードは自由に改変・転用できます。詳細は [LICENSE](LICENSE) を参照してください。
