# ch09-operations — 運用（監視・障害対応・スケール）

第9章「運用 ― 監視・障害対応・スケール（KRaft 前提）」の実機検証コード。

3 ブローカー KRaft クラスタを立て、レプリケーション・リーダー再選出・
コンシューマラグの監視を確認する。

## 構成

- `docker-compose.yml` — 3 ブローカー KRaft クラスタ（broker,controller 兼任）
- `examples/topic-rf3.sh` — replication.factor=3 + min.insync.replicas=2 のトピック作成と describe
- `examples/broker-down.sh` — ブローカー停止 → リーダー再選出 → 復帰 → ISR 回復
- `examples/lag.sh` — コンシューマラグを発生させて監視で捉える

## 実行

```bash
docker compose up -d
sleep 15   # 3 台のクラスタが揃うまで待つ

bash examples/topic-rf3.sh
bash examples/lag.sh
bash examples/broker-down.sh

docker compose down -v
```

## 確認できること

- `kafka-topics --describe` の Leader / Replicas / Isr が 3 台に分散している
- broker を 1 台止めると、その broker がリーダーだったパーティションのリーダーが
  ISR 内の別 broker に移り、止めた broker が Isr から消える
- broker を戻すと Isr が 1,2,3 に回復する（次の 1 台を止める前に待つべき状態）
- 消費を止めて produce を続けると、`kafka-consumer-groups --describe` の LAG が増える

検証環境: apache/kafka:4.3.1（3 ブローカー KRaft）/ Go 1.26
