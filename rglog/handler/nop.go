package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"sync"
)

// 何もしないハンドラ。
// デバッグ用。
type nopHandler struct {
	lock sync.Mutex
	lv   level.Level
}

func NewNopHandler() Handler {
	return &nopHandler{lv: level.ALL}
}

func (hndl *nopHandler) Level() level.Level {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	return hndl.lv
}

func (hndl *nopHandler) SetLevel(lv level.Level) {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	hndl.lv = lv
}

func (hndl *nopHandler) Output(rec Record) {}

func (hndl *nopHandler) Flush() {}

func (hndl *nopHandler) Close() {}
