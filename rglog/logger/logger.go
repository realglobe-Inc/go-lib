package logger

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
)

type Logger interface {
	// 登録されているハンドラ一覧を取得する。
	Handlers() []handler.Handler
	// ハンドラを追加する。
	AddHandler(hndl handler.Handler)
	// ハンドラを削除する。
	RemoveHandler(hndl handler.Handler)

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
