# ch07-transactions

第7章「配信保証とトランザクション」のサンプルコード。

- `main.go` — `orders` を read_committed で読み、value に `-processed` を付けて `orders-processed` にトランザクションで produce する consume-transform-produce（exactly-once）。`GroupTransactSession` で Begin → Produce → End(TryCommit) を回す。

## 前提

リポジトリのルートで Kafka を起動しておきます。

```bash
docker compose up -d
```

## 実行

```bash
# end-to-end（送信 → トランザクション処理 → ラグ確認 → read_committed で読む）
bash ch07-transactions/verify.sh
```

`main.go` 単体を動かす場合は、先に `orders` / `orders-processed` トピックを作り、`go run ./ch04-producer` で `orders` に送ってから `go run ./ch07-transactions` を実行します。デモのため、一定時間メッセージが来なければ自動で終了します。

## ポイント

- `kgo.TransactionalID(...)` でトランザクション ID を設定（ゾンビ Producer の締め出し）
- `kgo.FetchIsolationLevel(kgo.ReadCommitted())` でコミット済みだけ読む
- `GroupTransactSession` はリバランス時に commit を abort へ自動変換し、二重 produce を防ぐ
- 書き込みと「読んだオフセットのコミット」が `End` で原子的に確定する
