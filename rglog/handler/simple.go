package handler

import (
	"bufio"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type skelHandler struct {
	sync.Mutex
	lv level.Level

	formatter Formatter
	writer    io.Writer
}

func (hndl *skelHandler) Output(depth int, lv level.Level, v ...interface{}) {
	hndl.Lock()
	defer hndl.Unlock()

	// lv を読むためだけにロック。Go の仕様によっては必要無い。

	if lv > hndl.lv {
		return
	}

	// この辺は標準の log.Output を参考にした。
	// release lock while getting caller info - it's expensive.
	hndl.Unlock()
	now := time.Now()
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(depth + 1)
	if !ok {
		file = "???"
		line = 0
	}
	buff := hndl.formatter.Format(now, file, line, lv, v...)

	hndl.Lock()

	hndl.writer.Write(buff)
}

func (hndl *skelHandler) SetLevel(lv level.Level) {
	hndl.Lock()
	defer hndl.Unlock()

	hndl.lv = lv
}

func (hndl *skelHandler) Flush() {
}

func NewConsoleHandler() Handler {
	return NewConsoleHandlerUsing(simpleFormatter{})
}

func NewConsoleHandlerUsing(formatter Formatter) Handler {
	return &skelHandler{formatter: formatter, writer: os.Stderr}
}

type flushHandler struct {
	skelHandler
	writer *bufio.Writer
}

func (hndl *flushHandler) Flush() {
	hndl.Lock()
	defer hndl.Unlock()

	hndl.writer.Flush()
}

const logPerm = 0644

func NewFileHandler(path string) (Handler, error) {
	return NewFileHandlerUsing(path, simpleFormatter{})
}

func NewFileHandlerUsing(path string, formatter Formatter) (Handler, error) {
	// file の Close はプログラムの終処理任せ。
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, logPerm)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	writer := bufio.NewWriter(file)
	return &flushHandler{skelHandler{formatter: formatter, writer: writer}, writer}, nil
}
