package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"testing"
)

func TestSyslogHundler(t *testing.T) {
	// ただ使えるかだけ。

	hndl, err := NewSyslogHandler("go-lib-rg")
	if err != nil {
		t.Fatal(err)
	}
	hndl.SetLevel(level.ALL)
	hndl.Output(0, level.INFO, "test")
	hndl.Output(0, level.ERR, "test2")
	hndl.Flush()
}

// TODO 複数のコネクションで大量にログを吐くとデッドロックする場合がある。対処法不明。
func _TestSyslogHandler(t *testing.T) {
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
