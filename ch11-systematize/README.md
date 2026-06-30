# ch11-systematize — テスト・再現性

kfake(franz-go のインメモリ偽ブローカー)を使ったユニットテスト。
Docker を立てずに Producer/Consumer のロジックを高速に検証できる。

## 実行

```bash
go test -race ./ch11-systematize/...
```

外部プロセス(Docker など)は不要。kfake がインメモリで Kafka プロトコルを話す。

## テストの内容

- `produce_test.go` — `ControlKey(kmsg.Produce, ...)` で Produce リクエストを
  傍受し、送信した値が正しいかを検証する。
- `consume_test.go` — 偽クラスタへ実際に produce → consume するラウンドトリップ。

## 結合テスト

本物の Kafka に対する end-to-end は、ルートの `docker compose up -d` で
立てた Kafka 4.x KRaft に対して各章のコードを流す(Ch02〜Ch09 のディレクトリ)。
CI では Docker Kafka をサービスコンテナとして起動するか、testcontainers 的な
ライブラリでコンテナを起動する。

## バージョン

- franz-go v1.21.4 / kfake(franz-go の pkg/kfake)
