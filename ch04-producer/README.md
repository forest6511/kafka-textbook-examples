# ch04-producer

第4章「franz-go で Producer を書く」のサンプルコード。

`main.go` は本文の段階的なコード（最小の同期送信 → 非同期 → キー付き＋ヘッダー → graceful shutdown）を 1 本にまとめた実行サンプルです。

## 前提

リポジトリのルートで Kafka を起動しておきます。

```bash
docker compose up -d
```

`apache/kafka:4.3.1` を KRaft モード（ZooKeeper 不要）で起動します。

## 実行

```bash
# Producer を実行して CLI で確認まで通す
bash ch04-producer/verify.sh

# Producer だけ実行
go run ./ch04-producer
```

## 確認できること

- `ProduceSync` による同期送信と、`Produce` + コールバックによる非同期送信
- キー付き送信（`customer-42`）が単一パーティションに入ること（出力例: `partition=0 offset=2`。前に 2 件送っているためオフセットは 2）
- キーなし送信は uniform-bytes パーティショナで振り分けられること。少数を短時間に送ると 64KiB 単位のバッチにまとまり、同じパーティション（この例では 3 件とも partition 0）に入ること
- `Flush` でバッファを送り切ってから `Close` する graceful shutdown

実行環境: Apache Kafka 4.3.1（KRaft）/ franz-go v1.21.4 / Go 1.26
