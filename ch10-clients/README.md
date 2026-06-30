# ch10-clients — Go クライアントを書き比べる

franz-go / confluent-kafka-go / segmentio で同じ処理(Producer→Consumer)を書き、
ビルドと疎通の違いを実機で確認する。

## 実行

ルートで Kafka を起動してから:

```bash
docker compose up -d
bash ch10-clients/verify.sh
```

## 各クライアントの違い(実測)

- **franz-go** (`franz/`) — Pure Go。`CGO_ENABLED=0` でビルド・実行できる。
- **segmentio** (`segmentio/`) — Pure Go。`CGO_ENABLED=0` でビルド・実行できる。
- **confluent** (`confluent/`) — librdkafka(C)依存。`CGO_ENABLED=0` だと
  `undefined: kafka.ConfigMap` 等でビルド失敗する。`CGO_ENABLED=1`(+必要なら
  `-tags musl` / `-tags dynamic`)が要る。本リポジトリでは `confluent` ビルドタグで
  既定ビルドから除外している。

## バージョン(2026-07-01 時点)

- franz-go v1.21.4
- segmentio/kafka-go v0.4.51
- confluent-kafka-go v2.15.0

## 実機結果

```text
[franz-go] 3件を送信しました / 受信しました
[segmentio] 3件を送信しました / 受信しました
[confluent] 3件を送信しました / 受信しました
```

トピック ch10-compare には franz-0..2 / segmentio-0..2 / (confluent-0..2) が並ぶ。
各クライアントは別々のコンシューマグループなので、それぞれが先頭から読む。
