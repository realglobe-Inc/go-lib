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
	"container/list"
	"errors"
	"fmt"
	"github.com/realglobe-Inc/go-lib/erro"
	"os"
	"path/filepath"
	"strconv"
)

const (
	filePerm os.FileMode = 0644
	dirPerm              = 0755
)

// ファイルにログを書き込み、一杯になったらそのファイルをバックアップしてから
// 真っ新にしてまた書き込む coreHandler。
type rotateCoreHandler struct {
	// ログファイルのパス。
	path string
	// ログファイルを最大でいくつバックアップするか。
	num int
	// ログの書式。
	fmter Formatter

	// 以下、作業用データ。

	// 今開いているファイル。
	file *os.File
	// バッファ兼書き込み口。
	buff *fileLogBuffer
}

func (core *rotateCoreHandler) output(rec Record) {
	buff := core.fmter.Format(rec)

	const ( // てきとう。
		writeSize  = 4096
		bufferSize = 2*writeSize + 1024
	)

	core.buff.add(buff)
	defer core.buff.trim(bufferSize)
	core.flushCore(writeSize)
}

func (core *rotateCoreHandler) flush() {
	core.flushCore(0)
}

var fileIsFull error = errors.New("file is full")

func (core *rotateCoreHandler) flushCore(writeSize int) {
	// ロックファイルをつくったほうが良いが、OS 依存なので止めとく。

	for {
		if core.file == nil {
			// ファイルを開く。
			if err := os.MkdirAll(filepath.Dir(core.path), dirPerm); err != nil {
				fmt.Fprintln(os.Stderr, erro.Wrap(err))
				return
			}

			var err error
			core.file, err = os.OpenFile(core.path, os.O_RDWR|os.O_APPEND|os.O_CREATE, filePerm)
			if err != nil {
				fmt.Fprintln(os.Stderr, erro.Wrap(err))
				return
			}

			if err := core.buff.setSink(core.file); err != nil {
				fmt.Fprintln(os.Stderr, erro.Wrap(err))
				return
			}
		}

		// ファイルを開いている。

		err := core.buff.flushIfNeeded(writeSize)
		switch err {
		case nil:
			// 書き込み成功。
			return
		case fileIsFull:
			// ローテートする。
			core.file.Close()
			core.file = nil
			if err := rotateFile(core.path, core.num); err != nil {
				fmt.Fprintln(os.Stderr, erro.Wrap(err))
				return
			}
			continue
		default:
			fmt.Fprintln(os.Stderr, erro.Wrap(err))
			core.file.Close()
			core.file = nil
			return
		}
	}
}

func rotateFile(path string, num int) error {
	var n int

	for n = 1; n < num; n++ {
		_, err := os.Stat(path + "." + strconv.Itoa(n))
		if err == nil {
			continue
		} else if os.IsNotExist(err) {
			break
		} else {
			return erro.Wrap(err)
		}
	}
	n--

	// .{n} が残す中で一番最後。

	// .{i} を .{i+1} に。
	for ; n > 0; n-- {
		from := path + "." + strconv.Itoa(n)
		to := path + "." + strconv.Itoa(n+1)
		if err := os.Rename(from, to); err != nil {
			return erro.Wrap(err)
		}
	}

	// 最新版 を .1 に。
	if err := os.Rename(path, path+".1"); err != nil {
		return erro.Wrap(err)
	}

	return nil
}

func (core *rotateCoreHandler) close() {
	if core.file == nil {
		return
	}
	core.flush()
	core.buff.flush()
	core.file.Close()
	core.file = nil
}

func NewRotateHandler(path string, limit int64, num int) Handler {
	return NewRotateHandlerUsing(path, limit, num, SimpleFormatter)
}

func NewRotateHandlerUsing(path string, limit int64, num int, fmter Formatter) Handler {
	return wrapCoreHandler(newSynchronizedCoreHandler(&rotateCoreHandler{
		path:  path,
		num:   num,
		fmter: fmter,
		buff:  newFileLogBuffer(int(limit)),
	}))
}

// ファイル用のバッファ付き書き込みラッパー。
type fileLogBuffer struct {
	limit int // 32 ビット環境では 4GB 制限。

	// 書き込みファイル。
	sink     *os.File
	sinkSize int

	// まだ書き込んでない分。先頭が一番古いキュー。
	unwritten     *list.List
	unwrittenSize int
}

func newFileLogBuffer(limit int) *fileLogBuffer {
	return &fileLogBuffer{
		limit:     limit,
		unwritten: list.New(),
	}
}

func (buff *fileLogBuffer) setSink(sink *os.File) error {
	stat, err := sink.Stat()
	if err != nil {
		return erro.Wrap(err)
	}
	buff.sink = sink
	buff.sinkSize = int(stat.Size())
	return nil
}

func (buff *fileLogBuffer) add(log []byte) {
	buff.unwritten.PushBack(log)
	buff.unwrittenSize += len(log)
}

// サイズが baseline 未満にならないところまで古いログを消す。
func (buff *fileLogBuffer) trim(baseline int) {
	for buff.unwrittenSize > baseline && buff.unwritten.Len() > 0 {
		elem := buff.unwritten.Front()
		log := elem.Value.([]byte)
		if buff.unwrittenSize-len(log) < baseline {
			return
		}
		buff.unwritten.Remove(elem)
		buff.unwrittenSize -= len(log)
		fmt.Fprintln(os.Stderr, "Drop log: "+string(log))
	}
}

func (buff *fileLogBuffer) flushIfNeeded(unwrittenLimit int) error {
	if buff.unwrittenSize <= unwrittenLimit {
		return nil
	}
	return buff.flush()
}

func (buff *fileLogBuffer) flush() error {
	n := 0
	writeBuff := []byte{}
	for elem := buff.unwritten.Front(); elem != nil; elem = elem.Next() {
		log := elem.Value.([]byte)
		if curSize := buff.sinkSize + len(writeBuff); curSize > 0 && // サイズオーバーしても 1 個は必ず書く。
			curSize+len(log) > buff.limit {
			break
		}
		writeBuff = append(writeBuff, log...)
		n++
	}

	if _, err := buff.sink.Write(writeBuff); err != nil {
		return erro.Wrap(err)
	}

	buff.sinkSize += len(writeBuff)
	if n == buff.unwritten.Len() {
		// 全部書き出した。
		buff.unwritten.Init()
		buff.unwrittenSize = 0
		return nil
	}

	// n 個だけ書き出した。

	for i := 0; i < n; i++ {
		buff.unwritten.Remove(buff.unwritten.Front())
	}
	buff.unwrittenSize -= len(writeBuff)
	return fileIsFull
}

func (buff *fileLogBuffer) close() {
	buff.flush()
	buff.trim(0)
	*buff = fileLogBuffer{}
}
