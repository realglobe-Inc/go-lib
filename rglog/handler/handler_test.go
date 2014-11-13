package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"testing"
)

func testHandler(t *testing.T, hndl Handler) {
	// ただ走らせるだけ、確認できない。
	// panic で止まらないことの確認くらいにはなる。
	for _, lv := range level.Values() {
		hndl.SetLevel(lv)
		hndl.Output(1, lv, "test", lv)
	}

	hndl.Flush()
	hndl.Close()
}

func benchmarkHandler(b *testing.B, hndl Handler) {
	defer hndl.Close()
	hndl.SetLevel(level.ALL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hndl.Output(0, level.ERR, i)
	}
}
