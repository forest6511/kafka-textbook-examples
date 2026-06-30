// ch05-consumer: franz-go で Consumer を書く（章本文の最終形）
//
// 事前にルートで `docker compose up -d` を実行し、ch04-producer で
// orders にメッセージを送っておくこと。
//
//	go run ./ch04-producer   # 先にメッセージを送る
//	go run ./ch05-consumer   # 受け取る（Ctrl-C で停止）
//
// コンシューマグループ orders-processor に参加し、orders を最古から読む。
// 既定で 5 秒ごとに自動コミットし、Ctrl-C で安全に停止する。
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
		// グループに参加する。これで負荷分散とオフセット管理が効く。
		kgo.ConsumerGroup("orders-processor"),
		kgo.ConsumeTopics("orders"),
		// 未知のオフセットに出会ったら最古から読む（初回の全件読み用）。
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		log.Fatalf("クライアント生成に失敗: %v", err)
	}
	defer cl.Close()

	// Ctrl-C で ctx をキャンセルし、ポーリングループを安全に抜ける。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	fmt.Println("受信を開始します（Ctrl-C で停止）")
	for {
		fetches := cl.PollFetches(ctx)
		// ctx がキャンセルされたら抜ける（graceful shutdown）。
		if ctx.Err() != nil {
			fmt.Println("停止します")
			return
		}
		// 取得時のエラーは内部でリトライされるが、致命的なものはここに来る。
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Printf("取得エラー: %v", errs)
			continue
		}
		// 1 件ずつ処理する。既定で 5 秒ごとに自動コミットされる。
		fetches.EachRecord(func(r *kgo.Record) {
			fmt.Printf("受信: partition=%d offset=%d key=%s value=%s\n",
				r.Partition, r.Offset, string(r.Key), string(r.Value))
		})
	}
}
