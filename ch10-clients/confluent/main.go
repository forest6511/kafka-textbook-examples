//go:build confluent

// ch10-clients/confluent: confluent-kafka-go で同一処理を書く
//
// このファイルは `confluent` ビルドタグでのみコンパイルされる。
// librdkafka(C)が要るため、既定の `go build ./...` から除外している。
// ビルド/実行するには:
//
//	CGO_ENABLED=1 go run -tags confluent ./ch10-clients/confluent
//
// confluent-kafka-go は C ライブラリ librdkafka のラッパー。
// このため CGO_ENABLED=0 ではビルドできない(cgo 必須)。
//
//	CGO_ENABLED=1 go run ./ch10-clients/confluent
//
// Alpine(musl)では `-tags musl`、動的リンクは `-tags dynamic` を付ける。
// franz-go / segmentio が CGO_ENABLED=0 で通るのと対照的に、
// このディレクトリは C ツールチェインと librdkafka を要する。
package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	broker = "localhost:9092"
	topic  = "ch10-compare"
	group  = "ch10-confluent"
)

func main() {
	// Producer: Produce はチャネルに積み、配信レポートは Events() で受ける。
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatalf("プロデューサ生成に失敗: %v", err)
	}

	for i := 0; i < 3; i++ {
		value := fmt.Sprintf("confluent-%d", i)
		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: stringPtr(topic), Partition: kafka.PartitionAny},
			Value:          []byte(value),
		}, nil)
		if err != nil {
			log.Fatalf("送信に失敗: %v", err)
		}
	}
	// Flush は未送信を送り切るまで待つ(ミリ秒指定)。
	p.Flush(5000)
	p.Close()
	fmt.Println("[confluent] 3件を送信しました")

	// Consumer: コンシューマグループに参加して読む。
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("コンシューマ生成に失敗: %v", err)
	}
	defer c.Close()

	if err := c.Subscribe(topic, nil); err != nil {
		log.Fatalf("購読に失敗: %v", err)
	}

	for got := 0; got < 3; {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			log.Fatalf("受信に失敗: %v", err)
		}
		fmt.Printf("[confluent] 受信: %s\n", string(msg.Value))
		got++
	}
	fmt.Println("[confluent] 3件を受信しました")
}

func stringPtr(s string) *string { return &s }
