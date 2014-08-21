package handler

import (
	"bufio"
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// 与えられた出力先に書き出すだけの Handler.
// スレッドセーフ。
type basicHandler struct {
	sync.Mutex
	level.Level
	Formatter

	io.Writer
}

func NewBasicHandler(output io.Writer) Handler {
	return NewBasicHandlerUsing(output, SimpleFormatter)
}

func NewBasicHandlerUsing(output io.Writer, fmter Formatter) Handler {
	return newBasicHandlerUsing(output, fmter)
}

func newBasicHandlerUsing(output io.Writer, fmter Formatter) *basicHandler {
	return &basicHandler{Formatter: fmter, Writer: output}
}

func (hndl *basicHandler) Output(depth int, lv level.Level, v ...interface{}) {
	hndl.Lock()
	defer hndl.Unlock()

	if lv > hndl.Level {
		return
	}
	hndl.Unlock()

	// この辺は標準の log.Output を参考にした。
	// release lock while getting caller info - it's expensive.
	date := time.Now()
	_, file, line, ok := runtime.Caller(depth + 1)
	if ok {
		file = trimPrefix(file)
	} else {
		file = "???"
		line = 0
	}
	buff := hndl.Format(date, file, line, lv, v...)

	hndl.Lock()
	hndl.Write(buff)
}

func (hndl *basicHandler) SetLevel(lv level.Level) {
	hndl.Lock()
	defer hndl.Unlock()

	hndl.Level = lv
}

func (hndl *basicHandler) Flush() {
	return
}

func (hndl *basicHandler) Close() {
	return
}

// 与えられた出力先にバッファを挟んで書き出すだけの Handler.
// スレッドセーフ。
type flushHandler struct {
	*basicHandler
	*bufio.Writer
}

func (hndl *flushHandler) Flush() {
	hndl.Lock()
	defer hndl.Unlock()

	if err := hndl.Writer.Flush(); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
}

func NewFlushHandler(output io.Writer) Handler {
	return NewFlushHandlerUsing(output, SimpleFormatter)
}

func NewFlushHandlerUsing(output io.Writer, fmter Formatter) Handler {
	return newFlushHandlerUsing(output, fmter)
}

func newFlushHandlerUsing(output io.Writer, fmter Formatter) *flushHandler {
	bufOutput := bufio.NewWriter(output)
	return &flushHandler{newBasicHandlerUsing(bufOutput, fmter), bufOutput}
}

// Close にも対応。
// スレッドセーフ。
type closeHandler struct {
	*flushHandler
	io.Closer
}

func (hndl *closeHandler) Close() {
	hndl.Lock()
	defer hndl.Unlock()

	hndl.Flush()

	if err := hndl.Closer.Close(); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
}

func NewCloseHandler(output io.WriteCloser) Handler {
	return NewFlushHandlerUsing(output, SimpleFormatter)
}

func NewCloseHandlerUsing(output io.WriteCloser, fmter Formatter) Handler {
	return &closeHandler{newFlushHandlerUsing(output, fmter), output}
}

// 標準エラー出力に書き出す Handler。
func NewConsoleHandler() Handler {
	return NewConsoleHandlerUsing(SimpleFormatter)
}

func NewConsoleHandlerUsing(fmter Formatter) Handler {
	return NewBasicHandlerUsing(os.Stderr, fmter)
}

// 1 ファイルに延々と書き続ける Handler。
func NewFileHandler(path string) (Handler, error) {
	return NewFileHandlerUsing(path, SimpleFormatter)
}

func NewFileHandlerUsing(path string, fmter Formatter) (Handler, error) {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return nil, erro.Wrap(err)
	}
	// file の Close はプログラムの終処理任せ。
	output, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, filePerm)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return NewCloseHandlerUsing(output, fmter), nil
}
