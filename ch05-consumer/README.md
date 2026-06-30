# ch05-consumer

第5章「franz-go で Consumer を書く」のサンプルコード。

- `main.go` — コンシューマグループに参加し、`PollFetches` ループで読む基本形（既定の自動コミット + Ctrl-C で graceful shutdown）
- `manual-commit/main.go` — 自動コミットを切り、処理してからコミットする手動コミットの例（`OnPartitionsRevoked` でのコミット込み）

## 前提

リポジトリのルートで Kafka を起動し、`ch04-producer` でメッセージを送っておきます。

```bash
docker compose up -d
go run ./ch04-producer
```

## 実行

```bash
# end-to-end（送信 → 受信 → ラグ確認）を一気に通す
bash ch05-consumer/verify.sh

# Consumer だけ（Ctrl-C で停止）
go run ./ch05-consumer

# 手動コミット版
go run ./ch05-consumer/manual-commit
```

## 確認できること

- `kgo.ConsumerGroup` + `kgo.ConsumeTopics` でグループに参加し、パーティションが自動割り当てされること
- `PollFetches` でメッセージを取得し、`EachRecord` で 1 件ずつ処理すること
- 既定で 5 秒ごとに自動コミットされ、`kafka-consumer-groups --describe` の LAG が 0 になること
- `DisableAutoCommit` + `CommitUncommittedOffsets` + `OnPartitionsRevoked` による手動コミット

実行環境: Apache Kafka 4.3.1（KRaft）/ franz-go v1.21.4 / Go 1.26
