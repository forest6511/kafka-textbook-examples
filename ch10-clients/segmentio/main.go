// ch10-clients/segmentio: segmentio/kafka-go で同一処理を書く
//
// segmentio/kafka-go は Pure Go。高レベルの Writer(Producer)と
// Reader(Consumer)で同じ処理を書く。CGO_ENABLED=0 でビルドできる。
//
//	docker compose up -d                 # ルートで Kafka を起動
//	go run ./ch10-clients/segmentio      # produce 3件 → consume 3件
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	broker = "localhost:9092"
	topic  = "ch10-compare"
	group  = "ch10-segmentio"
)

func main() {
	ctx := context.Background()

	// Producer: Writer で 3件を送信する。RequireAll は全 ISR の ack を待つ。
	w := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}
	msgs := make([]kafka.Message, 3)
	for i := range msgs {
		msgs[i] = kafka.Message{Value: []byte(fmt.Sprintf("segmentio-%d", i))}
	}
	if err := w.WriteMessages(ctx, msgs...); err != nil {
		log.Fatalf("送信に失敗: %v", err)
	}
	w.Close()
	fmt.Println("[segmentio] 3件を送信しました")

	// Consumer: Reader に GroupID を渡すとコンシューマグループとして読む。
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: group,
		Topic:   topic,
	})
	defer r.Close()

	for got := 0; got < 3; got++ {
		// 全件読み終えても待ち続けないよう、タイムアウト付き ctx で読む。
		rctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		msg, err := r.ReadMessage(rctx)
		cancel()
		if err != nil {
			log.Fatalf("受信に失敗: %v", err)
		}
		fmt.Printf("[segmentio] 受信: %s\n", string(msg.Value))
	}
	fmt.Println("[segmentio] 3件を受信しました")
}
