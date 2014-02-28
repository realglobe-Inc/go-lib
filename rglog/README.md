rglog
==========

リアルグローブ式ロガー。

使い方
----------

必要なパッケージインポートする。
```
import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"github.com/realglobe-Inc/go-lib-rg/rglog"
)
```

プログラム実行のはじめのうちに標準動作を設定をする。
```
func main() {
	...

	// 使用する一番上の Logger を取得。
	log := rglog.GetLogger("github.com/realglobe-Inc/daiku")
	defer rglog.Flush()

	// これ以上は登らないことを明示。
	log.SetUseParent(false)
	// 表示できるレベルを設定。
	log.SetLevel(level.DEBUG)

	// INFO 以上を画面に表示する。
	hndl := handler.NewConsoleHandler()
	hndl.SetLevel(level.INFO)
	log.AddHandler(hndl)

	// 全てをファイルに出力する。
	hndl, err := handler.NewFileHandler(logPath)
	if err != nil {
		os.Exit(1)
	}
	hndl.SetLevel(level.DEBUG)
	log.AddHandler(hndl)

	...
}
```
一応、無設定だと "" の Logger まで遡って INFO 以上を画面に表示する。


使用したいところで、標準以下の Logger を取得して使う。
```
func Function() {
	...
	rglog.GetLogger("github.com/realglobe-Inc/daiku/change").Info("Logging message.")
	...
}
```

何度も使うなら init で取得しておくと良い。
```
var log rglog.Logger

func init() {
	log = rglog.GetLogger("github.com/realglobe-Inc/daiku/change")
}

func Function() {
	...
	log.Info("Logging message.")
	...
}
```

標準とは違う動作や追加動作をさせたい場合は、取得したロガーに AddHandler やら SetLevel やら SetUseParent やらで。
