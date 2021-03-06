// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog/level"
)

// 与えられた出力先に書き出すだけの Handler.
// スレッドセーフ。
type basicHandler struct {
	lock  sync.Mutex
	lv    level.Level
	fmter Formatter

	sink io.Writer
}

func NewBasicHandler(sink io.Writer) Handler {
	return NewBasicHandlerUsing(sink, SimpleFormatter)
}

func NewBasicHandlerUsing(sink io.Writer, fmter Formatter) Handler {
	return newBasicHandlerUsing(sink, fmter)
}

func newBasicHandlerUsing(sink io.Writer, fmter Formatter) *basicHandler {
	return &basicHandler{lv: level.ALL, fmter: fmter, sink: sink}
}

func (hndl *basicHandler) Level() level.Level {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	return hndl.lv
}

func (hndl *basicHandler) SetLevel(lv level.Level) {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	hndl.lv = lv
}

func (hndl *basicHandler) Output(rec Record) {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	if !rec.Level().Lower(hndl.lv) {
		hndl.sink.Write(hndl.fmter.Format(rec))
	}
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
	flusher *bufio.Writer
}

func (hndl *flushHandler) Flush() {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	if err := hndl.flusher.Flush(); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
}

func NewFlushHandler(sink io.Writer) Handler {
	return NewFlushHandlerUsing(sink, SimpleFormatter)
}

func NewFlushHandlerUsing(sink io.Writer, fmter Formatter) Handler {
	return newFlushHandlerUsing(sink, fmter)
}

func newFlushHandlerUsing(sink io.Writer, fmter Formatter) *flushHandler {
	bufSink := bufio.NewWriter(sink)
	return &flushHandler{newBasicHandlerUsing(bufSink, fmter), bufSink}
}

// Close にも対応。
// スレッドセーフ。
type closeHandler struct {
	*flushHandler
	closer io.Closer
}

func (hndl *closeHandler) Close() {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	hndl.Flush()

	if err := hndl.closer.Close(); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
}

func NewCloseHandler(sink io.WriteCloser) Handler {
	return NewFlushHandlerUsing(sink, SimpleFormatter)
}

func NewCloseHandlerUsing(sink io.WriteCloser, fmter Formatter) Handler {
	return &closeHandler{newFlushHandlerUsing(sink, fmter), sink}
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
	sink, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, filePerm)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return NewCloseHandlerUsing(sink, fmter), nil
}
