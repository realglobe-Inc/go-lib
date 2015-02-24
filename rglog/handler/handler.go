package handler

import (
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"time"
)

// ログの書き出し機。
type Handler interface {
	// 書き出すレベル。
	// Output の引数の重要度が Handler の重要度以上であれば実際に書き出し、未満であれば書き出さない。
	// 初期値は基本的に level.ALL。
	Level() level.Level
	// 書き出すレベルを指定する。
	SetLevel(lv level.Level)

	// 書き出す。
	Output(rec Record)

	// バッファを使っているなら、低層に書き出す。
	Flush()

	Close()
}

type Record interface {
	Date() time.Time
	Level() level.Level
	File() string
	Line() int
	Message() string
}
