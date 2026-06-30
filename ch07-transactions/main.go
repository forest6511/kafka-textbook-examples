// ch07-transactions: consume-transform-produce をトランザクションで原子的に行う例。
//
// orders を読み、value に "-processed" を付けて orders-processed に produce し、
// 読んだオフセットのコミットと produce を 1 つのトランザクションでまとめて確定する。
// GroupTransactSession がリバランス時の安全（commit→abort 自動切替）を担う。
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	sess, err := kgo.NewGroupTransactSession(
		kgo.SeedBrokers("localhost:9092"),
		kgo.TransactionalID("orders-eos-processor"),
		kgo.ConsumerGroup("orders-eos"),
		kgo.ConsumeTopics("orders"),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		kgo.FetchIsolationLevel(kgo.ReadCommitted()),
		kgo.DefaultProduceTopic("orders-processed"),
	)
	if err != nil {
		log.Fatalf("セッション生成に失敗: %v", err)
	}
	defer sess.Close()

	ctx := context.Background()
	processed := 0
	// デモのため、一定時間メッセージが来なければ終了する。
	idle := 0
	for {
		fetches := sess.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Fatalf("取得エラー: %v", errs)
		}
		n := fetches.NumRecords()
		if n == 0 {
			idle++
			if idle >= 3 {
				break
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}
		idle = 0

		if err := sess.Begin(); err != nil {
			log.Fatalf("トランザクション開始に失敗: %v", err)
		}

		fetches.EachRecord(func(r *kgo.Record) {
			out := &kgo.Record{
				Key:   r.Key,
				Value: []byte(string(r.Value) + "-processed"),
			}
			sess.Produce(ctx, out, func(_ *kgo.Record, err error) {
				if err != nil {
					log.Printf("produce 失敗: %v", err)
				}
			})
			processed++
		})

		committed, err := sess.End(ctx, kgo.TryCommit)
		if err != nil {
			log.Fatalf("トランザクション終了に失敗: %v", err)
		}
		if committed {
			fmt.Printf("コミット: %d 件を orders-processed へ\n", n)
		} else {
			fmt.Printf("アボート: リバランス等のため %d 件を取り消し\n", n)
		}
	}
	fmt.Printf("終了: 合計 %d 件を処理しました\n", processed)
}
