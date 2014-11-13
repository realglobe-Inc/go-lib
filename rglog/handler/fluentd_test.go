package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"net"
	"testing"
)

// テストしたかったら fluentd サーバーを立ててから。
var fluentdAddr = "localhost:24224"

func init() {
	if fluentdAddr != "" {
		// 実際にサーバーが立っているかどうか調べる。
		// 立ってなかったらテストはスキップ。
		conn, err := net.Dial("tcp", fluentdAddr)
		if err != nil {
			fluentdAddr = ""
		} else {
			conn.Close()
		}
	}
}

// ただ使えるかだけ。
func TestFluentdHandler(t *testing.T) {
	if fluentdAddr == "" {
		t.SkipNow()
	}

	hndl := NewFluentdHandler(fluentdAddr, "rglog.test")
	defer hndl.Flush()

	hndl.SetLevel(level.ALL)
	hndl.Output(0, level.INFO, "test")
	hndl.Output(0, level.ERR, "test2")
}

// 色んな長さのメッセージを送る。
// MessagePack 部分のテスト。
func TestFluentdHandlerMessageLength(t *testing.T) {
	if fluentdAddr == "" {
		t.SkipNow()
	}

	hndl := NewFluentdHandler(fluentdAddr, "rglog.test")
	defer hndl.Flush()

	hndl.SetLevel(level.ALL)

	for i := 0; i < 18; i++ {
		msg := ""
		for j := 0; j < (1 << uint(i)); j++ {
			msg += "a"
		}
		hndl.Output(0, level.INFO, msg)
	}
}

// 複数接続で。
func TestManyFluentdHandler(t *testing.T) {
	if fluentdAddr == "" {
		t.SkipNow()
	}

	n := 20
	loop := 100

	hndls := []Handler{}
	for i := 0; i < n; i++ {
		hndl := NewFluentdHandler(fluentdAddr, "rglog.test")
		defer hndl.Flush()

		hndl.SetLevel(level.ALL)
		hndls = append(hndls, hndl)
	}

	for i := 0; i < loop; i++ {
		for j := 0; j < len(hndls); j++ {
			hndls[j].Output(0, level.ERR, "a ", j, i)
		}
	}
}

func BenchmarkFluentdHandler(b *testing.B) {
	if fluentdAddr == "" {
		b.SkipNow()
	}

	hndl := NewFluentdHandler(fluentdAddr, "rglog.test")
	benchmarkHandler(b, hndl)
}
