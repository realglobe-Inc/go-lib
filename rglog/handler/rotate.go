package handler

import (
	"bufio"
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"os"
	"path/filepath"
	"strconv"
	"time"
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
	// ログファイル 1 つの最大サイズ。
	limit int64
	// ログファイルを最大でいくつバックアップするか。
	num int
	// ログの書式。
	fmter Formatter

	// 以下、作業用データ。

	// 今開いているファイル。
	file *os.File
	// 今のサイズ。
	size int64
	// 書き込み口。
	sink *bufio.Writer
}

func (hndl *rotateCoreHandler) output(file string, line int, lv level.Level, v ...interface{}) {
	if err := hndl.outputCore(file, line, lv, v...); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
}

func (hndl *rotateCoreHandler) outputCore(file string, line int, lv level.Level, v ...interface{}) error {
	// ロックファイルをつくったほうが良いが、OS 依存なので止めとく。

	buff := hndl.fmter.Format(time.Now(), file, line, lv, v...)

	for {
		if hndl.file == nil {
			// ファイルを開く。
			if err := os.MkdirAll(filepath.Dir(hndl.path), dirPerm); err != nil {
				return erro.Wrap(err)
			}

			var err error
			hndl.file, err = os.OpenFile(hndl.path, os.O_RDWR|os.O_APPEND|os.O_CREATE, filePerm)
			if err != nil {
				return erro.Wrap(err)
			}

			stat, err := hndl.file.Stat()
			if err != nil {
				return erro.Wrap(err)
			}
			hndl.size = stat.Size()
		}

		// ファイルを開いている。

		if hndl.size != 0 && hndl.size+int64(len(buff)) > hndl.limit { // ファイルは必ず 1 度は使う。
			// ローテートする。
			if err := hndl.rotate(); err != nil {
				return erro.Wrap(err)
			}
			continue
		}

		if hndl.sink == nil {
			hndl.sink = bufio.NewWriter(hndl.file)
		}

		// ファイルに余裕がある。
		break
	}

	size, err := hndl.sink.Write(buff)
	if err != nil {
		return erro.Wrap(err)
	}

	hndl.size += int64(size)
	return nil
}

func (hndl *rotateCoreHandler) rotate() error {
	if hndl.file != nil {
		if hndl.sink != nil {
			if err := hndl.sink.Flush(); err != nil {
				return erro.Wrap(err)
			}
			hndl.sink = nil
		}
		if err := hndl.file.Close(); err != nil {
			return erro.Wrap(err)
		}
		hndl.file = nil
		hndl.size = 0
	}

	return erro.Wrap(rotateFile(hndl.path, hndl.num))
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

func (hndl *rotateCoreHandler) flush() {
	if hndl.sink == nil {
		return
	}
	if err := hndl.sink.Flush(); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
}

func (hndl *rotateCoreHandler) close() {
	hndl.flush()

	if hndl.file == nil {
		return
	}
	if err := hndl.file.Close(); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
}

func NewRotateHandler(path string, limit int64, num int) Handler {
	return NewRotateHandlerUsing(path, limit, num, SimpleFormatter)
}

func NewRotateHandlerUsing(path string, limit int64, num int, fmter Formatter) Handler {
	return wrapCoreHandler(newSynchronizedCoreHandler(&rotateCoreHandler{path: path, limit: limit, num: num, fmter: fmter}))
}
