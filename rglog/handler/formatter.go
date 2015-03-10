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
	"fmt"
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"strconv"
)

type Formatter interface {
	Format(rec Record) []byte
}

// {日時} {レベル} {ファイル名}:{行番号} {メッセージ}
type simpleFormatter struct{}

var SimpleFormatter = &simpleFormatter{}

// Level の最短幅。
var lvWidth int

func init() {
	lvWidth = 10
	for _, lv := range level.Values() {
		if w := len(lv.String()); w < lvWidth {
			lvWidth = w
		}
	}
}

func (formatter simpleFormatter) Format(rec Record) []byte {
	year, month, day := rec.Date().Date()
	hour, min, sec := rec.Date().Clock()
	microSec := rec.Date().Nanosecond() / 1000

	buff := fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d.%06d %."+strconv.Itoa(lvWidth)+"v %s:%d %s\n",
		year, int(month), day, hour, min, sec, microSec, rec.Level(), rec.File(), rec.Line(), rec.Message())

	return []byte(buff)
}

// [{レベル}] {メッセージ}
type levelOnlyFormatter struct{}

var LevelOnlyFormatter = &levelOnlyFormatter{}

func (formatter levelOnlyFormatter) Format(rec Record) []byte {
	buff := fmt.Sprintf("[%."+strconv.Itoa(lvWidth)+"v] %s\n", rec.Level(), rec.Message())
	return []byte(buff)
}
