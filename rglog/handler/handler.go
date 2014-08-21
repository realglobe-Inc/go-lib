package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
)

type Handler interface {
	SetLevel(lv level.Level)

	Output(depth int, lv level.Level, v ...interface{}) // depth は Logger が重ねたスタックの数。
	Flush()
	Close()
}
