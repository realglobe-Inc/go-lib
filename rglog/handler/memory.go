package handler

import (
	"bytes"
)

// デバッグ用のハンドラ。
type MemoryHandler struct {
	*basicHandler
	buff *bytes.Buffer
}

func NewMemoryHandler() *MemoryHandler {
	return NewMemoryHandlerUsing(SimpleFormatter)
}

func NewMemoryHandlerUsing(fmter Formatter) *MemoryHandler {
	return newMemoryHandlerUsing(fmter)
}

func newMemoryHandlerUsing(fmter Formatter) *MemoryHandler {
	var buff bytes.Buffer
	return &MemoryHandler{newBasicHandlerUsing(&buff, fmter), &buff}
}

func (hndl *MemoryHandler) Dump() string {
	return hndl.buff.String()
}
