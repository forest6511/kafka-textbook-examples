# ch02-setup — Docker で Kafka を動かし、franz-go から最小疎通する

第2章のサンプルです。KRaft モードの Kafka をローカルに起動し、franz-go から接続できることを確認します。

## 1. Kafka を起動する

リポジトリのルートで実行します（共通の `docker-compose.yml` を使います）。

```bash
docker compose up -d
docker compose ps
```

`STATUS` が `Up` になっていれば起動完了です。`localhost:9092` で待ち受けます。

## 2. CLI で疎通を確認する

トピック一覧を取得します。まだ何も作っていないので、空が返れば正常です。

```bash
docker compose exec broker /opt/kafka/bin/kafka-topics.sh \
    --bootstrap-server localhost:9092 --list
```

## 3. franz-go から接続する

```bash
go run ./ch02-setup
```

期待する出力:

```text
Kafka に接続できました (localhost:9092)
```

`main.go` は `kgo.NewClient` でクライアントを作り、`cl.Ping(ctx)` でブローカーへの到達性だけを確認します。トピックを作らずに「繋がるか」だけを見たいときの最小コードです。

## 4. 後片付け

```bash
docker compose down -v      # 停止 + データ削除
```

## 動作確認環境

- Apache Kafka 4.x（`apache/kafka` 公式イメージ・KRaft モード）
- franz-go v1.21+ / Go 1.26
