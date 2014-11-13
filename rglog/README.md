rglog
==========

リアルグローブ式ロガー。

使い方
----------

プログラム実行のはじめのうちに標準動作を設定する。

```
import (
	...
	"github.com/realglobe-Inc/go-lib-rg/rglog"
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	...
)

...

func main() {
	...
	// 使用する一番上の Logger を取得。
	log := rglog.Logger("github.com/realglobe-Inc")
	defer rglog.Flush()

	// これより上には遡らない。
	log.SetUseParent(false)
	// 書き出すレベルを設定。
	log.SetLevel(level.ALL)

	// 画面に INFO 以上を表示する。
	hndl := handler.NewConsoleHandler()
	hndl.SetLevel(level.INFO)
	log.AddHandler("console", hndl)

	// fileSize バイト、backupNum 個までのファイル logPath にデバッグ情報を出力する。
	hndl = handler.NewRotateHandler(logPath, fileSize, backupNum)
	hndl.SetLevel(level.DEBUG)
	log.AddHandler("file", hndl)
	...
}
```

無設定だと "" の Logger まで遡り、"" の Logger には "console" という名前で INFO 以上を画面に表示するハンドラが登録されている。


使用したいところで、適当な Logger を取得して使う。

```
import (
	...
	"github.com/realglobe-Inc/go-lib-rg/rglog"
	...
)

...

func Function() {
	...
	rglog.Logger("github.com/realglobe-Inc/a/b/c").Info("Log message")
	...
}
```

何度も使うなら初期化時に取得しておくと良い。

```
var log = rglog.Logger("github.com/realglobe-Inc/a/b/c")

...

func Function() {
	...
	log.Info("Log message")
	...
}
```

標準とは異なる動作を追加したいときは、log.AddHandler やら log.SetLevel やらで。
log.SetUseParent(false) すれば、標準動作をさせないようにできる。

ログメッセージの生成が重い場合、IsLoggable でログを取らない場合には飛ばすようにできる。

```
func Function() {
	...
	if log.IsLoggable(level.INFO) {
		// ログの書き出しが発生する場合のみ、ここに至る。
		log.Info(generateLogMessage())
	}
	...
}
```
