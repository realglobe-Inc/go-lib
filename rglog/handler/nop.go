package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
)

// 何もしないハンドラ。
// デバッグ用。
type nopHandler struct{}

func NewNopHandler() Handler {
	return &nopHandler{}
}

func (hndl *nopHandler) SetLevel(lv level.Level) {}

func (hndl *nopHandler) Output(depth int, lv level.Level, v ...interface{}) {}

func (hndl *nopHandler) Flush() {}

func (hndl *nopHandler) Close() {}
