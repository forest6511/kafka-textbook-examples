// ch08-schema: Confluent Schema Registry を franz-go の sr パッケージで操作する例。
//
// 1. User v1 スキーマを登録し、振られた ID を確認する
// 2. ID からスキーマを引き直す
// 3. subject を BACKWARD 互換に設定する
// 4. 互換な変更(任意フィールド追加)が通ることを確認する
// 5. 壊す変更(必須フィールド追加)が弾かれることを確認する
// 6. sr.Serde でワイヤフォーマット(マジックバイト + 4 バイト ID + 本体)を符号化・復号する
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/twmb/franz-go/pkg/sr"
)

// User は本文で扱う最小の値。Serde では Avro を組まず、
// JSON エンコードを通じてワイヤフォーマット(ID プレフィックス)を見せる。
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func mustRead(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("スキーマ読み込みに失敗: %v", err)
	}
	return string(b)
}

func main() {
	ctx := context.Background()

	cl, err := sr.NewClient(sr.URLs("http://localhost:8081"))
	if err != nil {
		log.Fatalf("レジストリ接続に失敗: %v", err)
	}

	const subject = "user-value"

	// 1. v1 を登録する。戻り値に subject・version・グローバル ID が入る。
	v1 := sr.Schema{Schema: mustRead("examples/user-v1.json"), Type: sr.TypeAvro}
	ss, err := cl.CreateSchema(ctx, subject, v1)
	if err != nil {
		log.Fatalf("v1 登録に失敗: %v", err)
	}
	fmt.Printf("登録: subject=%s version=%d id=%d\n", ss.Subject, ss.Version, ss.ID)

	// 2. ID からスキーマを引き直す。
	got, err := cl.SchemaByID(ctx, ss.ID)
	if err != nil {
		log.Fatalf("ID から取得に失敗: %v", err)
	}
	fmt.Printf("ID %d から取得した型: %s\n", ss.ID, got.Type)

	// 3. subject を BACKWARD 互換に設定する(既定値だが明示する)。
	for _, r := range cl.SetCompatibility(ctx,
		sr.SetCompatibility{Level: sr.CompatBackward}, subject) {
		if r.Err != nil {
			log.Fatalf("互換性設定に失敗: %v", r.Err)
		}
		fmt.Printf("互換性モード: %s\n", r.Level)
	}

	// 4. 互換な変更(default 付きフィールド email を追加)をチェックする。
	v2ok := sr.Schema{Schema: mustRead("examples/user-v2-ok.json"), Type: sr.TypeAvro}
	okRes, err := cl.CheckCompatibility(ctx, subject, ss.Version, v2ok)
	if err != nil {
		log.Fatalf("互換性チェックに失敗: %v", err)
	}
	fmt.Printf("email 追加(default あり)は互換か: %v\n", okRes.Is)

	// 5. 壊す変更(default なしの必須フィールド age を追加)をチェックする。
	v2bad := sr.Schema{Schema: mustRead("examples/user-v2-bad.json"), Type: sr.TypeAvro}
	badRes, err := cl.CheckCompatibility(ctx, subject, ss.Version, v2bad)
	if err != nil {
		log.Fatalf("互換性チェックに失敗: %v", err)
	}
	fmt.Printf("age 追加(default なし)は互換か: %v\n", badRes.Is)
	if !badRes.Is && len(badRes.Messages) > 0 {
		fmt.Printf("弾かれた理由: %s\n", badRes.Messages[0])
	}

	// 6. sr.Serde でワイヤフォーマットを符号化・復号する。
	//    Register で「ID」と「Go の値・エンコード/デコード関数」を結びつける。
	serde := sr.NewSerde()
	serde.Register(ss.ID, User{},
		sr.EncodeFn(func(v any) ([]byte, error) { return json.Marshal(v) }),
		sr.DecodeFn(func(b []byte, v any) error { return json.Unmarshal(b, v) }),
	)

	encoded, err := serde.Encode(User{ID: "u-1", Name: "Alice"})
	if err != nil {
		log.Fatalf("符号化に失敗: %v", err)
	}
	fmt.Printf("符号化後の先頭 5 バイト: % x\n", encoded[:5])

	id, body, err := serde.DecodeID(encoded)
	if err != nil {
		log.Fatalf("ID 取り出しに失敗: %v", err)
	}
	fmt.Printf("先頭から読んだスキーマ ID: %d / 本体: %s\n", id, string(body))

	var back User
	if err := serde.Decode(encoded, &back); err != nil {
		log.Fatalf("復号に失敗: %v", err)
	}
	fmt.Printf("復号した値: %+v\n", back)
}
