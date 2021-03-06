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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/realglobe-Inc/go-lib/rglog/handler"
	"github.com/realglobe-Inc/go-lib/rglog/level"
)

func testLoggerHandler(t *testing.T, mgr Manager) {
	log := mgr.Logger("a/b/c/d")

	if hndl := log.Handler("test"); hndl != nil {
		t.Fatal(hndl)
	}

	memHndl := handler.NewNopHandler()
	if hndl := log.AddHandler("test", memHndl); hndl != nil {
		t.Fatal(hndl)
	}

	oldHndl := log.Handler("test")
	if oldHndl != memHndl {
		t.Fatal(oldHndl, memHndl)
	}

	memHndl = handler.NewNopHandler()
	if hndl := log.AddHandler("test", memHndl); hndl != oldHndl {
		t.Fatal(hndl, oldHndl)
	}

	if hndl := log.RemoveHandler("test"); hndl != memHndl {
		t.Fatal(hndl, memHndl)
	}

	if hndl := log.Handler("test"); hndl != nil {
		t.Fatal(hndl)
	}
}

func testLoggerLevel(t *testing.T, mgr Manager) {
	log := mgr.Logger("a/b/c/d")

	if log.Level() != level.OFF {
		t.Fatal(log.Level())
	}

	for _, lv := range level.Values() {
		log.SetLevel(lv)
		if log.Level() != lv {
			t.Fatal(log.Level(), lv)
		}
	}
}

func testLoggerUseParent(t *testing.T, mgr Manager) {
	log := mgr.Logger("a/b/c/d")

	for _, b := range []bool{true, false} {
		log.SetUseParent(b)
		if log.UseParent() != b {
			t.Fatal(log.UseParent(), b)
		}
	}
}

func testLoggerIsLoggable(t *testing.T, mgr Manager) {
	//                    先祖にハンドラが無い 先祖の基準重要度より低い 先祖の基準重要度より高い
	// ハンドラが無い     false                false                    true
	// 基準重要度より低い false                false                    true
	// 基準重要度より高い true                 true                     true

	log := mgr.Logger("a/b/c/d")
	log.SetLevel(level.INFO)
	log.SetUseParent(true)
	parentLog := mgr.Logger("a/b")
	parentLog.SetLevel(level.INFO)
	parentLog.SetUseParent(false)

	// ハンドラが無い、先祖にハンドラが無い。
	if log.IsLoggable(level.INFO) {
		t.Fatal("true: no handler, no handler")
	}

	log.AddHandler("test", handler.NewNopHandler())

	// 基準重要度より低い、先祖にハンドラが無い。
	if log.IsLoggable(level.DEBUG) {
		t.Fatal("true: lower level, no handler")
	}

	// 基準重要度より高い、先祖にハンドラが無い。
	if !log.IsLoggable(level.INFO) {
		t.Fatal("false: upper or equal level, no handler")
	}

	log.RemoveHandler("test")
	parentLog.AddHandler("test", handler.NewNopHandler())
	parentLog.SetLevel(level.WARN)

	// ハンドラが無い、先祖の基準重要度より低い。
	if log.IsLoggable(level.INFO) {
		t.Fatal("true: no handler, lower level")
	}

	log.AddHandler("test", handler.NewNopHandler())

	// 基準重要度より低い、先祖の基準重要度より低い。
	if log.IsLoggable(level.DEBUG) {
		t.Fatal("true: lower level, lower level")
	}

	// 基準重要度より高い、先祖の基準重要度より低い。
	if !log.IsLoggable(level.INFO) {
		t.Fatal("false: upper or equal level, lower level")
	}

	log.RemoveHandler("test")
	parentLog.AddHandler("test", handler.NewNopHandler())
	parentLog.SetLevel(level.DEBUG)

	// ハンドラが無い、先祖の基準重要度より高い。
	if !log.IsLoggable(level.INFO) {
		t.Fatal("false: no handler, upper level")
	}

	log.AddHandler("test", handler.NewNopHandler())

	// 基準重要度より低い、先祖の基準重要度より高い。
	if !log.IsLoggable(level.DEBUG) {
		t.Fatal("true: lower level, upper level")
	}

	// 基準重要度より高い、先祖の基準重要度より高い。
	if !log.IsLoggable(level.INFO) {
		t.Fatal("false: upper or equal level, upper level")
	}
}

func testLoggerLog(t *testing.T, mgr Manager) {
	// ただやってみるだけ。
	// すぐに panic にならないことの確認くらいにはなる。

	log := mgr.Logger("a/b/c/d")

	log.SetLevel(level.INFO)

	for _, lv := range level.Values() {
		log.Log(lv, "test message level ", lv)
	}
}

// ファイル名がちゃんと解決できてるかどうか。
func testLoggerFileName(t *testing.T, mgr Manager) {
	log := mgr.Logger("a/b/c")
	log.SetLevel(level.ALL)
	log.SetUseParent(false)

	hndl := handler.NewMemoryHandlerUsing(&fileOnlyFormatter{})
	hndl.SetLevel(level.ALL)
	log.AddHandler("test", hndl)

	log.Log(level.INFO, "")
	dum := hndl.Dump()
	if dum != filepath.Join("github.com", "realglobe-Inc", "go-lib", "rglog", "logger", "logger_test.go") {
		t.Fatal(dum)
	}

	for _, logging := range []func(...interface{}){log.Err, log.Warn, log.Info, log.Debug} {
		hndl := handler.NewMemoryHandlerUsing(&fileOnlyFormatter{})
		hndl.SetLevel(level.ALL)
		log.AddHandler("test", hndl)

		logging("")
		dum := hndl.Dump()
		if dum != filepath.Join("github.com", "realglobe-Inc", "go-lib", "rglog", "logger", "logger_test.go") {
			t.Fatal(dum)
		}
	}
}

type fileOnlyFormatter struct{}

func (fmter *fileOnlyFormatter) Format(rec handler.Record) []byte {
	return []byte(rec.File())
}

func testLoggerConcurrent(t *testing.T, mgr Manager) {
	rootLabel := "a/b/c"

	conc := 100
	loop := 1000

	rootLog := mgr.Logger(rootLabel)
	rootLog.SetLevel(level.ALL)
	rootLog.SetUseParent(false)

	path := filepath.Join(os.TempDir(), "locklog_test")
	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}
	defer os.Remove(path)

	hndl := handler.NewMemoryHandler()
	hndl.SetLevel(level.DEBUG)
	rootLog.AddHandler("test", hndl)

	c := make(chan bool)
	timeout := time.After(time.Duration(int64(conc*loop*100) * int64(time.Microsecond)))

	for i := 0; i < conc; i++ {
		id := i
		go func() {
			for j := 0; j < loop; j++ {
				mgr.Logger(rootLabel+"/"+strconv.Itoa(id)).Info(id, j)
			}

			c <- true
		}()
	}

	for i := 0; i < conc; i++ {
		select {
		case <-c:
		case <-timeout:
			t.Fatal("Dead lock?")
		}
	}

	mgr.Flush()

	// ファイルに書き込めているかどうか検査。
	buff := hndl.Dump()

	if len(buff) > 0 && buff[len(buff)-1] == '\n' { // 最後の空行は抜かしとく。
		buff = buff[:len(buff)-1]
	}

	lines := strings.Split(buff, "\n")
	if len(lines) != conc*loop {
		t.Fatal(len(lines), conc*loop)
	}
}

func benchmarkLoggerConcurrent(b *testing.B, mgr Manager) {
	rootLabel := "a/b/c"
	conc := 100

	rootLog := mgr.Logger(rootLabel)
	rootLog.SetLevel(level.ALL)
	rootLog.SetUseParent(false)

	path := filepath.Join(os.TempDir(), "locklog_test")
	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			b.Fatal(err)
		}
	}
	defer os.Remove(path)

	hndl := handler.NewNopHandler()
	hndl.SetLevel(level.DEBUG)
	rootLog.AddHandler("test", hndl)

	c := make(chan bool)
	timeout := time.After(time.Duration(int64(conc*b.N*100) * int64(time.Microsecond)))

	for i := 0; i < conc; i++ {
		id := i
		go func() {
			for j := 0; j < b.N/conc; j++ {
				mgr.Logger(rootLabel+"/"+strconv.Itoa(id)).Info(id, j)
			}

			c <- true
		}()
	}

	for i := 0; i < conc; i++ {
		select {
		case <-c:
		case <-timeout:
			b.Fatal("Dead lock?")
		}
	}

	mgr.Flush()
}
