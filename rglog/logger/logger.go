package logger

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
)

type Logger interface {
	// 登録してあるハンドラを取得する。
	Handler(key string) handler.Handler
	// ハンドラを登録する。
	// 既に同じ key でハンドラが登録してあったら、新しい方に置き換えて、古い方を返す。
	AddHandler(key string, hndl handler.Handler) (oldHndl handler.Handler)
	// ハンドラを登録から外す。
	RemoveHandler(key string) (oldHndl handler.Handler)

	// ハンドラに処理させる重要度の下限を返す。
	Level() level.Level
	// ハンドラに処理させる重要度の下限を指定する。
	SetLevel(lv level.Level)

	// GetLogger に渡す name が / で区切られた木構造を表すとして、親の Logger にも処理させるかどうか。
	SetUseParent(useParent bool)

	Err(v ...interface{})
	Warn(v ...interface{})
	Info(v ...interface{})
	Debug(v ...interface{})
}
