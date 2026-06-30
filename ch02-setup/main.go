// ch02-setup: franz-go で Kafka に最小疎通する。
//
// docker compose up -d で起動した KRaft クラスタ(localhost:9092)に接続し、
// ブローカーのメタデータを取得して、接続できることだけを確認する。
//
//	go run ./ch02-setup
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	// SeedBrokers は最初に問い合わせるブローカー。ここからクラスタ全体を発見する。
	cl, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
	)
	if err != nil {
		log.Fatalf("クライアント生成に失敗: %v", err)
	}
	defer cl.Close()

	// Ping はブローカーへの到達性だけを確認する。トピックを作らずに疎通を見たいときに使う。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cl.Ping(ctx); err != nil {
		log.Fatalf("ブローカーに到達できません: %v", err)
	}

	fmt.Println("Kafka に接続できました (localhost:9092)")
}
