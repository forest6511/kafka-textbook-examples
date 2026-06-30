// ch10-clients/franz: franz-go で同一処理(Producer→Consumer)を書く
//
// 3クライアント(franz-go / confluent-kafka-go / segmentio)で同じ処理を
// 書き比べるための franz-go 版。Pure Go なので CGO_ENABLED=0 でビルドできる。
//
//	docker compose up -d            # ルートで Kafka 4.x KRaft を起動
//	go run ./ch10-clients/franz     # produce 3件 → consume 3件
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	broker = "localhost:9092"
	topic  = "ch10-compare"
	group  = "ch10-franz"
)

func main() {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(broker),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		log.Fatalf("クライアント生成に失敗: %v", err)
	}
	defer cl.Close()

	ctx := context.Background()

	// Producer: 3件を同期送信する。
	for i := 0; i < 3; i++ {
		rec := &kgo.Record{Topic: topic, Value: []byte(fmt.Sprintf("franz-%d", i))}
		if err := cl.ProduceSync(ctx, rec).FirstErr(); err != nil {
			log.Fatalf("送信に失敗: %v", err)
		}
	}
	fmt.Println("[franz-go] 3件を送信しました")

	// Consumer: 3件を受信する。
	got := 0
	for got < 3 {
		fetches := cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Fatalf("取得に失敗: %v", errs)
		}
		fetches.EachRecord(func(r *kgo.Record) {
			fmt.Printf("[franz-go] 受信: %s\n", string(r.Value))
			got++
		})
	}
	fmt.Println("[franz-go] 3件を受信しました")
}
