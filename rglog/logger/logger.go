// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"github.com/realglobe-Inc/go-lib/rglog/handler"
	"github.com/realglobe-Inc/go-lib/rglog/level"
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
	// 初期値は level.OFF。
	Level() level.Level
	// ハンドラに処理させる重要度の下限を指定する。
	SetLevel(lv level.Level)

	// 識別子を / 区切りの木構造として、親の識別子のロガーにも処理させるかどうか。
	UseParent() bool
	// 識別子を / 区切りの木構造として、親の識別子のロガーにも処理させるかどうかを指定する。
	SetUseParent(useParent bool)

	// 指定した重要度でログを取ったら、ハンドラに処理させるかどうか。
	// UseParent が true な限りの先祖ロガーも含む。
	IsLoggable(lv level.Level) bool

	// ログを取る。
	Log(lv level.Level, v ...interface{})
	// Log(level.ERR, v...) と一緒。
	Err(v ...interface{})
	// Log(level.WARN, v...) と一緒。
	Warn(v ...interface{})
	// Log(level.INFO, v...) と一緒。
	Info(v ...interface{})
	// Log(level.DEBUG, v...) と一緒。
	Debug(v ...interface{})
}

type Manager interface {
	Logger(name string) Logger
	Flush()
}
