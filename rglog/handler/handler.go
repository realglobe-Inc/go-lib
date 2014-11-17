package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"time"
)

// ログの書き出し機。
type Handler interface {
	// 書き出すレベルを指定する。
	// Output において、ここで指定した以上の重要度であれば実際に書き出し、未満の重要度であれば書き出さない。
	// 初期値は基本的に level.ALL。
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
