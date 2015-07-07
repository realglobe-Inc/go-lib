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

package handler

import (
	"time"

	"github.com/realglobe-Inc/go-lib/rglog/level"
)

// ログの書き出し機。
type Handler interface {
	// 書き出すレベル。
	// Output の引数の重要度が Handler の重要度以上であれば実際に書き出し、未満であれば書き出さない。
	// 初期値は基本的に level.ALL。
	Level() level.Level
	// 書き出すレベルを指定する。
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
