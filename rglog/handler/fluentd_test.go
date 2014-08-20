package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"testing"
)

// テストしたかったら fluentd サーバーを立ててから。
var fluentdAddr = "localhost:24224"

// ただ使えるかだけ。
func _TestFluentdHundler(t *testing.T) {
	hndl, err := NewFluentdHandler(fluentdAddr, "rglog.test")
	if err != nil {
		t.Fatal(err)
	}
	defer hndl.Flush()

	hndl.SetLevel(level.ALL)
	hndl.Output(0, level.INFO, "test")
	hndl.Output(0, level.ERR, "test2")
}

// 色んな長さのメッセージを送る。
// MessagePack 部分のテスト。
func _TestFluentdHundlerMessageLength(t *testing.T) {
	hndl, err := NewFluentdHandler(fluentdAddr, "rglog.test")
	if err != nil {
		t.Fatal(err)
	}
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
func _TestManyFluentdHandler(t *testing.T) {
	n := 20
	loop := 100

	hndls := []Handler{}
	for i := 0; i < n; i++ {
		hndl, err := NewFluentdHandler(fluentdAddr, "rglog.test")
		if err != nil {
			t.Fatal(err)
		}
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

func _BenchmarkFluentdHandler(b *testing.B) {
	hndl, err := NewFluentdHandler(fluentdAddr, "rglog.test")
	if err != nil {
		b.Fatal(err)
	}
	benchmarkHandler(b, hndl)
}
