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
