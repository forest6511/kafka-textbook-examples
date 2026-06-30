// ch11-systematize: kfake を使った Producer のユニットテスト
//
// kfake はインメモリの偽 Kafka クラスタ。外部プロセス(Docker など)を
// 立てずに Producer/Consumer のロジックを高速に検証できる。
//
//	go test -race ./ch11-systematize/...
package main

import (
	"context"
	"testing"

	"github.com/twmb/franz-go/pkg/kfake"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
)

// sendOrder は「アプリ側の送信ロジック」を模した関数。実コードでは
// ここにビジネスロジックが入る。テスト対象はこの関数。
func sendOrder(ctx context.Context, cl *kgo.Client, topic, value string) error {
	rec := &kgo.Record{Topic: topic, Value: []byte(value)}
	return cl.ProduceSync(ctx, rec).FirstErr()
}

func TestSendOrder(t *testing.T) {
	const topic = "orders"

	// インメモリの偽クラスタを起動する。Docker は不要。
	cluster, err := kfake.NewCluster(kfake.SeedTopics(1, topic))
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Close()

	cl, err := kgo.NewClient(kgo.SeedBrokers(cluster.ListenAddrs()...))
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()

	// ControlKey で Produce リクエストを傍受し、送信内容を検証する。
	var gotValue string
	cluster.ControlKey(int16(kmsg.Produce), func(r kmsg.Request) (kmsg.Response, error, bool) {
		req := r.(*kmsg.ProduceRequest)
		b := kmsg.NewRecordBatch()
		if err := b.ReadFrom(req.Topics[0].Partitions[0].Records); err != nil {
			t.Error(err)
		}
		rr := kmsg.NewRecord()
		if err := rr.ReadFrom(b.Records); err != nil {
			t.Error(err)
		}
		gotValue = string(rr.Value)
		// false を返すと、偽クラスタが通常どおり応答を返す。
		return nil, nil, false
	})

	if err := sendOrder(context.Background(), cl, topic, "order-42"); err != nil {
		t.Fatalf("送信に失敗: %v", err)
	}
	if gotValue != "order-42" {
		t.Fatalf("送信値が一致しません: got %q, want %q", gotValue, "order-42")
	}
}
