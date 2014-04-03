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
