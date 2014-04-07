package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"testing"
)

func _TestFluentdHundler(t *testing.T) {
	// ただ使えるかだけ。

	hndl, err := NewFluentdHandler("localhost:8888", "oykt")
	if err != nil {
		t.Fatal(err)
	}
	hndl.SetLevel(level.ALL)
	hndl.Output(0, level.INFO, "test")
	hndl.Output(0, level.ERR, "test2")
	hndl.Flush()
}

func _TestFluentdHandler(t *testing.T) {
	n := 20
	loop := 100

	hndls := []Handler{}
	for i := 0; i < n; i++ {
		hndl, err := NewFluentdHandler("localhost:8888", "oykt")
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

func _BenchmarkFluentdHandler(b *testing.B) {
	hndl, err := NewFluentdHandler("localhost:8888", "oykt")
	if err != nil {
		b.Fatal(err)
	}
	hndl.SetLevel(level.ALL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hndl.Output(0, level.ERR, i)
	}
}
