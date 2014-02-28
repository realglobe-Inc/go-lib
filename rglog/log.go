package rglog

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"github.com/realglobe-Inc/go-lib-rg/rglog/locklog"
)

// Loggerのインターフェース
type Logger interface {
	AddHandler(hndl handler.Handler)
	RemoveHandler(hndl handler.Handler)

	SetLevel(lv level.Level)
	SetUseParent(use bool) // GetLogger に渡す name が / で区切られた木構造を表すとして、親にも処理させるかどうか。

	Err(v ...interface{})
	Warn(v ...interface{})
	Info(v ...interface{})
	Debug(v ...interface{})
}

// 各パッケージの init で 1 回だけ呼ぶくらいを想定。
func GetLogger(name string) Logger {
	return locklog.Logger(name)
}

// TODO 名前は Close でも良いが。
// TODO 手動で Flush しなくちゃならないのは面倒。終処理にフックしたい。
func Flush() {
	locklog.Flush()
}
