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
	"log/syslog"
	"testing"
	"time"
)

// 実際テストしたかったら true に。
var testSyslogHandlerFlag = true

func init() {
	if testSyslogHandlerFlag {
		// 実際にサーバーが立っているかどうか調べる。
		// 立ってなかったらテストはスキップ。
		conn, err := syslog.New(syslog.LOG_INFO, "test")
		if err != nil {
			testSyslogHandlerFlag = false
		} else {
			conn.Close()
		}
	}
}

func TestSyslogHandlerLevel(t *testing.T) {
	if !testSyslogHandlerFlag {
		t.SkipNow()
	}

	testHandlerLevel(t, NewSyslogHandler("go-lib"))
}

func TestSyslogHandlerOutput(t *testing.T) {
	if !testSyslogHandlerFlag {
		t.SkipNow()
	}

	testHandlerOutput(t, NewSyslogHandler("go-lib"))
}

// TODO 複数のコネクションで大量にログを吐くとデッドロックする場合がある。対処法不明。
func TestManySyslogHandler(t *testing.T) {
	if !testSyslogHandlerFlag {
		t.SkipNow()
	}

	n := 20
	loop := 100

	hndls := []Handler{}
	for i := 0; i < n; i++ {
		hndl := NewSyslogHandler("a")
		defer hndl.Close()
		hndls = append(hndls, hndl)
	}

	for i := 0; i < loop; i++ {
		for j := 0; j < len(hndls); j++ {
			hndls[j].Output(&record{date: time.Now(), lv: level.ERR, msg: fmt.Sprint("a ", j, i)})
		}
	}
}

func BenchmarkSyslogHandler(b *testing.B) {
	if !testSyslogHandlerFlag {
		b.SkipNow()
	}

	benchmarkHandler(b, NewSyslogHandler("go-lib"))
}
