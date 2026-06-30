// ch11-systematize: kfake を使った produce→consume ラウンドトリップのテスト
//
// ControlKey で傍受せず、偽クラスタに実際に produce して consume できるかを
// 確認する。アプリ側の「送って受け取る」一連を高速に検証できる。
package main

import (
	"context"
	"testing"
	"time"

	"github.com/twmb/franz-go/pkg/kfake"
	"github.com/twmb/franz-go/pkg/kgo"
)

func TestRoundTrip(t *testing.T) {
	const topic = "events"

	cluster, err := kfake.NewCluster(kfake.SeedTopics(1, topic))
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Close()

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(cluster.ListenAddrs()...),
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()

	ctx := context.Background()
	if err := cl.ProduceSync(ctx, &kgo.Record{Topic: topic, Value: []byte("hello")}).FirstErr(); err != nil {
		t.Fatalf("送信に失敗: %v", err)
	}

	// タイムアウト付きで 1 件受信する。
	pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	fetches := cl.PollFetches(pctx)
	if errs := fetches.Errors(); len(errs) > 0 {
		t.Fatalf("取得に失敗: %v", errs)
	}

	var got string
	fetches.EachRecord(func(r *kgo.Record) { got = string(r.Value) })
	if got != "hello" {
		t.Fatalf("受信値が一致しません: got %q, want %q", got, "hello")
	}
}
