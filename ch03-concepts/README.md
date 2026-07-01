# ch03-concepts — 中核概念を CLI で観察する

第3章のサンプルです。トピック・パーティション・オフセット・コンシューマグループを、
コードではなく CLI で観察します。事前にルートで `docker compose up -d` を実行しておきます。

## 観察 1: パーティションごとのオフセットとラグ

`groups.sh` は `orders` トピックに 12 件のキー付きメッセージを投入し、
`ch03-demo` グループで全件読んだあとに、グループの読み進み具合を表示します。

```bash
bash ch03-concepts/groups.sh
```

`kafka-consumer-groups.sh --describe` の出力で、パーティション 0/1/2 それぞれの
`CURRENT-OFFSET` / `LOG-END-OFFSET` / `LAG` が確認できます。読み切っていれば `LAG` は 0 です。

## 観察 2: 同じキーは同じパーティションへ

キー付きで投入したメッセージが、キーのハッシュで決まるパーティションに振り分けられ、
パーティションごとに件数が偏ることを `--describe` の `LOG-END-OFFSET` から読み取れます。

## 後片付け

```bash
docker compose down -v
```

## 動作確認環境

- Apache Kafka 4.3.1（apache/kafka 公式イメージ・KRaft モード）/ Go 1.26
