// ch05-consumer/manual-commit: 手動コミットの最小例。
//
// 自動コミットを切り、バッチを処理し切ってからコミットする。
// リバランスでのデータロス/重複を避けるため、OnPartitionsRevoked でもコミットする。
//
//	go run ./ch05-consumer/manual-commit   # Ctrl-C で停止
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
		kgo.ConsumerGroup("orders-manual"),
		kgo.ConsumeTopics("orders"),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		// 自動コミットを切る。コミットの責任を自分で持つ。
		kgo.DisableAutoCommit(),
		// 再割り当て前に未コミット分をコミットして取りこぼしを防ぐ。
		kgo.OnPartitionsRevoked(
			func(ctx context.Context, cl *kgo.Client, _ map[string][]int32) {
				if err := cl.CommitUncommittedOffsets(ctx); err != nil {
					log.Printf("revoke 時コミットに失敗: %v", err)
				}
			}),
	)
	if err != nil {
		log.Fatalf("クライアント生成に失敗: %v", err)
	}
	defer cl.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for {
		fetches := cl.PollFetches(ctx)
		if ctx.Err() != nil {
			return
		}
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Printf("取得エラー: %v", errs)
			continue
		}
		fetches.EachRecord(func(r *kgo.Record) {
			fmt.Printf("処理: offset=%d value=%s\n", r.Offset, r.Value)
		})
		// このバッチを処理し切ってからコミットする（処理 → コミットの順）。
		if err := cl.CommitUncommittedOffsets(ctx); err != nil {
			log.Printf("コミットに失敗: %v", err)
		}
	}
}
