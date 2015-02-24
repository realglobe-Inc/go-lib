package rglog

import (
	"github.com/realglobe-Inc/go-lib/rglog/handler"
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"github.com/realglobe-Inc/go-lib/rglog/logger"
)

// level ログの重要度。
// handler 書き出し機。
// logger ハンドラをまとめたり、親子関係をつくったり。

var mgr logger.Manager

func init() {
	mgr = logger.NewLockLoggerManager()
}

// 無設定時の動作を設定。
func init() {
	log := mgr.Logger("")
	log.SetLevel(level.INFO)
	log.SetUseParent(false)

	hndl := handler.NewConsoleHandler()
	hndl.SetLevel(level.ALL)
	log.AddHandler("console", hndl)
}

// 各パッケージの init で 1 回だけ呼ぶくらいを想定。
func Logger(name string) logger.Logger {
	return mgr.Logger(name)
}

// TODO 手動で Flush しなくちゃならないのは面倒。終処理にフックしたい。
func Flush() {
	mgr.Flush()
}
