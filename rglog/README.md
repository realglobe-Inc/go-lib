<!--
Copyright 2015 realglobe, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->


# rglog

リアルグローブ式ロガー。


## 1. 使い方

プログラム実行のはじめのうちに標準動作を設定する。

```Go
import (
	...
	"github.com/realglobe-Inc/go-lib/rglog"
	"github.com/realglobe-Inc/go-lib/rglog/handler"
	"github.com/realglobe-Inc/go-lib/rglog/level"
	...
)

...

func main() {
	...
	// 使用する一番上の logger.Logger を取得する。
	// 引数には import パス等を使用する。
	log := rglog.Logger("a/b/c")
	defer rglog.Flush()
	// これより上には遡らせない。
	log.SetUseParent(false)
	// 全てのログを登録されている handler.Handler に渡させる。logger.Logger の初期レベルは基本的に level.OFF。
	log.SetLevel(level.ALL)

	// 標準エラー出力に level.INFO 以上を書き出させる。
	hndl := handler.NewConsoleHandler()
	hndl.SetLevel(level.INFO)
	log.AddHandler("console", hndl)

	// ファイル path に、最大 size バイト、最大 n ファイルで、デバッグ情報まで書き出させる。
	hndl = handler.NewRotateHandler(path, size, n)
	// handler.Handler の初期レベルは基本的に level.ALL。
	log.AddHandler("file", hndl)
	...
}
```

無設定だと "" の logger.Logger まで遡る。
"" の logger.Logger は level.INFO 以上を登録されている handler.Handler に渡し、渡された全てのログを標準エラー出力に書き出す handler.Handler が "console" という名前で登録されている。
つまり、無設定だと level.INFO 以上が標準エラー出力に書き出される。

使用したいところで、適当な Logger を取得して使う。

```Go
import (
	...
	"github.com/realglobe-Inc/go-lib/rglog"
	...
)

...

func Function() {
	...
	rglog.Logger("a/b/c/d").Info("Log message")
	...
}
```

何度も使うなら初期化時に取得しておくと良い。

```Go
var log = rglog.Logger("a/b/c/d")

...

func Function() {
	...
	log.Info("Log message")
	...
}
```

標準とは異なる動作を追加したいときは、好きな handler.Handler を log.AddHandler したり、log.SetLevel したりで。
標準動作をさせないなら、log.SetUseParent(false)。

ログメッセージの生成が重いなら、IsLoggable で handler.Handler にログが渡されない場合には飛ばすこともできる。

```Go
func Function() {
	...
	if log.IsLoggable(level.INFO) {
		// 1 つ以上の Handler.handler にログが渡される場合、ここに至る。
		log.Log(level.INFO, generateLogMessage())
	}
	...
}
```


## 2. API

[GoDoc](http://godoc.org/github.com/realglobe-Inc/go-lib/rglog)


## 3. ライセンス

Apache License, Version 2.0
