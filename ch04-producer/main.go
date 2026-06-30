// ch04-producer: franz-go で Producer を書く（章本文の最終形）
//
// 事前にルートで `docker compose up -d` を実行しておくこと。
// 本書第4章の段階的なコードを 1 本にまとめた実行サンプル。
//
//	go run ./ch04-producer
//
// 送信後は次のコマンドで届いているか確認できる:
//
//	docker compose exec broker /opt/kafka/bin/kafka-console-consumer.sh \
//	    --bootstrap-server localhost:9092 --topic orders \
//	    --from-beginning --timeout-ms 3000 --property print.key=true
package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	// クライアントを生成する。Producer も Consumer も同じ *kgo.Client。
	// franz-go は既定でべき等プロデューサ（idempotent）で送るため、
	// enable.idempotence のような設定を明示しなくても重複が防がれる。
	cl, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
		// 既定は AllISRAcks（全 in-sync replica の書き込みを待つ）。
		// ここでは既定を明示している。弱めると idempotent と両立しない。
		kgo.RequiredAcks(kgo.AllISRAcks()),
		// 圧縮はバッチ単位。優先順で渡すとブローカー対応のものが選ばれる。
		kgo.ProducerBatchCompression(kgo.Lz4Compression(), kgo.NoCompression()),
		// 低スループットではわずかに待ってバッチ化すると効率が上がる。
		kgo.ProducerLinger(5*time.Millisecond),
	)
	if err != nil {
		log.Fatalf("クライアント生成に失敗: %v", err)
	}
	defer cl.Close()

	ctx := context.Background()

	// 1) 最小の同期送信。ProduceSync はバッチをまとめて送り、最初のエラーを返す。
	rec := &kgo.Record{Topic: "orders", Value: []byte("order-1 created")}
	if err := cl.ProduceSync(ctx, rec).FirstErr(); err != nil {
		log.Fatalf("同期送信に失敗: %v", err)
	}
	fmt.Println("同期送信 OK: order-1")

	// 2) 非同期送信。コールバックで完了を受け取り、WaitGroup で待ち合わせる。
	var wg sync.WaitGroup
	wg.Add(1)
	cl.Produce(ctx, &kgo.Record{
		Topic: "orders",
		Value: []byte("order-2 created"),
	}, func(_ *kgo.Record, err error) {
		defer wg.Done()
		if err != nil {
			log.Printf("非同期送信でエラー: %v", err)
			return
		}
		fmt.Println("非同期送信 OK: order-2")
	})
	wg.Wait()

	// 3) キー付き送信。同じキーは同じパーティションに入り、順序が保たれる。
	//    さらにヘッダーでメタデータ（イベント種別）を値と分けて運ぶ。
	keyed := &kgo.Record{
		Topic: "orders",
		Key:   []byte("customer-42"),
		Value: []byte("order-3 created"),
		Headers: []kgo.RecordHeader{
			{Key: "event-type", Value: []byte("OrderCreated")},
		},
	}
	r, err := cl.ProduceSync(ctx, keyed).First()
	if err != nil {
		log.Fatalf("キー付き送信に失敗: %v", err)
	}
	fmt.Printf("キー付き送信 OK: partition=%d offset=%d\n", r.Partition, r.Offset)

	// 4) graceful shutdown: バッファを送り切ってから閉じる。
	//    defer cl.Close() の前に Flush することで取りこぼしを防ぐ。
	if err := cl.Flush(ctx); err != nil {
		log.Fatalf("Flush に失敗: %v", err)
	}
	fmt.Println("全送信を Flush しました")
}
