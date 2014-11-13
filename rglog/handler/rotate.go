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

func (core *rotateCoreHandler) output(file string, line int, lv level.Level, v ...interface{}) {
	if err := core.outputCore(file, line, lv, v...); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
}

func (core *rotateCoreHandler) outputCore(file string, line int, lv level.Level, v ...interface{}) error {
	// ロックファイルをつくったほうが良いが、OS 依存なので止めとく。

	buff := core.fmter.Format(time.Now(), file, line, lv, v...)

	for {
		if core.file == nil {
			// ファイルを開く。
			if err := os.MkdirAll(filepath.Dir(core.path), dirPerm); err != nil {
				return erro.Wrap(err)
			}

			var err error
			core.file, err = os.OpenFile(core.path, os.O_RDWR|os.O_APPEND|os.O_CREATE, filePerm)
			if err != nil {
				return erro.Wrap(err)
			}

			stat, err := core.file.Stat()
			if err != nil {
				return erro.Wrap(err)
			}
			core.size = stat.Size()
		}

		// ファイルを開いている。

		if core.size != 0 && core.size+int64(len(buff)) > core.limit { // ファイルは必ず 1 度は使う。
			// ローテートする。
			if err := core.rotate(); err != nil {
				return erro.Wrap(err)
			}
			continue
		}

		if core.sink == nil {
			core.sink = bufio.NewWriter(core.file)
		}

		// ファイルに余裕がある。
		break
	}

	size, err := core.sink.Write(buff)
	if err != nil {
		return erro.Wrap(err)
	}

	core.size += int64(size)
	return nil
}

func (core *rotateCoreHandler) rotate() error {
	if core.file != nil {
		if core.sink != nil {
			if err := core.sink.Flush(); err != nil {
				return erro.Wrap(err)
			}
			core.sink = nil
		}
		if err := core.file.Close(); err != nil {
			return erro.Wrap(err)
		}
		core.file = nil
		core.size = 0
	}

	return erro.Wrap(rotateFile(core.path, core.num))
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

func (core *rotateCoreHandler) flush() {
	if core.sink == nil {
		return
	}
	if err := core.sink.Flush(); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
}

func (core *rotateCoreHandler) close() {
	core.flush()

	if core.file == nil {
		return
	}
	if err := core.file.Close(); err != nil {
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
