package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"testing"
)

func benchmarkHandler(b *testing.B, hndl Handler) {
	defer hndl.Flush()
	hndl.SetLevel(level.ALL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hndl.Output(0, level.ERR, i)
	}
}
