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
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"strconv"
	"testing"
	"time"
)

func testHandlerLevel(t *testing.T, hndl Handler) {
	if hndl.Level() != level.ALL {
		t.Fatal(hndl.Level())
	}

	for _, lv := range level.Values() {
		hndl.SetLevel(lv)
		if hndl.Level() != lv {
			t.Fatal(hndl.Level(), lv)
		}
	}
}

func testHandlerOutput(t *testing.T, hndl Handler) {
	// ただやってみるだけ。
	// すぐに panic にならないことの確認くらいにはなる。

	hndl.SetLevel(level.INFO)

	for _, lv := range level.Values() {
		hndl.Output(&record{time.Now(), lv, "test", 0, lv.String()})
	}

	hndl.Flush()
	hndl.Close()
}

func benchmarkHandler(b *testing.B, hndl Handler) {
	defer hndl.Close()
	hndl.SetLevel(level.ALL)

	b.ResetTimer()
	date := time.Now()
	for i := 0; i < b.N; i++ {
		hndl.Output(&record{date.Add(time.Duration(i) * time.Nanosecond), level.INFO, "test", 0, strconv.Itoa(i)})
	}
}

type record struct {
	date time.Time
	lv   level.Level
	file string
	line int
	msg  string
}

func (rec *record) Date() time.Time {
	return rec.date
}
func (rec *record) Level() level.Level {
	return rec.lv
}
func (rec *record) File() string {
	return rec.file
}
func (rec *record) Line() int {
	return rec.line
}
func (rec *record) Message() string {
	return rec.msg
}
