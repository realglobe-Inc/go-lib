package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"log/syslog"
	"testing"
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

func TestSyslogHandler(t *testing.T) {
	if !testSyslogHandlerFlag {
		t.SkipNow()
	}

	hndl, err := NewSyslogHandler("go-lib-rg")
	if err != nil {
		t.Fatal(err)
	}
	testHandler(t, hndl)
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
		hndl, err := NewSyslogHandler("a")
		if err != nil {
			t.Fatal(err)
		}
		hndl.SetLevel(level.ALL)
		hndls = append(hndls, hndl)
	}

	for i := 0; i < loop; i++ {
		for j := 0; j < len(hndls); j++ {
			hndls[j].Output(0, level.ERR, "a ", j, i)
		}
	}
}

func BenchmarkSyslogHandler(b *testing.B) {
	if !testSyslogHandlerFlag {
		b.SkipNow()
	}

	hndl, err := NewSyslogHandler("go-lib-rg")
	if err != nil {
		b.Fatal(err)
	}
	hndl.SetLevel(level.ALL)
	benchmarkHandler(b, hndl)
}
