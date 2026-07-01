# ch08-schema — スキーマと Schema Registry

第8章「スキーマとデータ整合 ― Schema Registry」の実機検証コード。

Confluent Schema Registry を franz-go の `sr` パッケージで操作し、
スキーマ登録・ID 参照・互換性チェック・ワイヤフォーマット符号化を確認する。

## 構成

- `docker-compose.yml` — Kafka 4.x(KRaft)+ Confluent Schema Registry(8081)
- `main.go` — sr.Client / sr.Serde の操作例
- `examples/user-v1.json` — 初版スキーマ(id, name)
- `examples/user-v2-ok.json` — 互換な変更(default 付き email を追加)
- `examples/user-v2-bad.json` — 壊す変更(default なしの age を追加)

## 実行

> **この章は `ch08-schema/` に `cd` してから実行します。** 章専用の `docker-compose.yml`（Kafka + Schema Registry）を使うため、ルートで Kafka を起動している場合は先に `docker compose down -v` で止めておいてください（ポート衝突を防ぐため）。

```bash
cd ch08-schema           # この章のディレクトリへ移動

docker compose up -d

# Schema Registry が立ち上がるまで待つ(8081 が 200 を返すまで)
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8081/subjects

go run .

docker compose down -v
```

## 期待される出力

```text
登録: subject=user-value version=1 id=1
ID 1 から取得した型: AVRO
互換性モード: BACKWARD
email 追加(default あり)は互換か: true
age 追加(default なし)は互換か: false
符号化後の先頭 5 バイト: 00 00 00 00 01
先頭から読んだスキーマ ID: 1 / 本体: {"id":"u-1","name":"Alice"}
復号した値: {ID:u-1 Name:Alice}
```

先頭 5 バイト `00 00 00 00 01` は、マジックバイト `00` + 4 バイトのスキーマ ID(=1)。
この ID で「どのスキーマか」を運ぶのが Schema Registry 運用の中核。

## REST で確認する

```bash
curl -s http://localhost:8081/subjects
# ["user-value"]

curl -s http://localhost:8081/subjects/user-value/versions
# [1]

curl -s http://localhost:8081/config/user-value
# {"compatibilityLevel":"BACKWARD"}
```

検証環境: apache/kafka:4.3.1 / confluentinc/cp-schema-registry:7.9.0 /
franz-go v1.21.4 / sr v1.7.0 / Go 1.26
