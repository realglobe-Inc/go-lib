package handler

import (
	"bufio"
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type FormatEntry struct {
	Date time.Time
	File string
	Line int
	Lv   level.Level
	Args []interface{}
}

// てきとうなポインタを Flush のトリガにする。
// ポインタ自体で判断するので、中身を保護する必要は無い。
var flushTrigger = &FormatEntry{}

func FlushTrigger() *FormatEntry {
	return flushTrigger
}

type goHandler struct {
	sync.Mutex
	lv level.Level

	done  chan bool
	queue chan *FormatEntry
}

func (hndl *goHandler) Output(depth int, lv level.Level, v ...interface{}) {
	hndl.Lock()
	if lv > hndl.lv {
		return
	}
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

	hndl.queue <- &FormatEntry{now, file, line, lv, v}
}

func (hndl *goHandler) SetLevel(lv level.Level) {
	hndl.Lock()
	defer hndl.Unlock()

	hndl.lv = lv
}

func (hndl *goHandler) Flush() {
	hndl.queue <- FlushTrigger()
	<-hndl.done
}

func NewGoHandler(goFunc func(done chan bool, queue chan *FormatEntry), queueCapacity int) Handler {
	done := make(chan bool)
	queue := make(chan *FormatEntry, queueCapacity)
	go goFunc(done, queue)
	return &goHandler{done: done, queue: queue}
}

const defaultQueueCapacity = 8192

func NewRotateHandler(path string, limit int64, num int) Handler {
	return NewRotateHandlerUsing(path, limit, num, defaultQueueCapacity, &simpleFormatter{})
}

const errThreshold = 5                    // 再試行する限度。
const coolDownDuration = time.Millisecond // 異常発生時に空ける間隔。

func NewRotateHandlerUsing(path string, limit int64, num, queueCapacity int, formatter Formatter) Handler {
	return NewGoHandler(func(done chan bool, queue chan *FormatEntry) {
		for errCount := 0; errCount < errThreshold; {

			// 異常が発生してたら、ちょっと落ち着く。
			if errCount > 0 {
				time.Sleep(coolDownDuration)
			}

			// ファイルを開く。
			file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, logPerm)
			if err != nil {
				fmt.Fprintln(os.Stderr, erro.Wrap(err))
				errCount++
				break
			}

			fi, err := file.Stat()
			if err != nil {
				fmt.Fprintln(os.Stderr, erro.Wrap(err))
				errCount++
				file.Close()
				break
			}

			size := fi.Size()
			writer := bufio.NewWriter(file)

			// ログを取る。
			for size <= limit { // 最大 1 エントリ分はみ出す。でも、limit == 0 でも動く。
				ent := <-queue

				if ent == FlushTrigger() {
					err := writer.Flush()
					done <- (err == nil) // デッドロック防止のため、先に返す。
					if err != nil {
						fmt.Fprintln(os.Stderr, erro.Wrap(err))
						break
					}
				}

				buff := formatter.Format(ent.Date, ent.File, ent.Line, ent.Lv, ent.Args...)
				_, err := writer.Write(buff)
				if err != nil {
					fmt.Fprintln(os.Stderr, erro.Wrap(err))
					break
				}

				size += int64(len(buff))
			}

			// ファイルを閉じる。
			if e := writer.Flush(); e != nil {
				fmt.Fprintln(os.Stderr, erro.Wrap(e))
				errCount++
				file.Close()
				break
			}

			if e := file.Close(); e != nil {
				fmt.Fprintln(os.Stderr, erro.Wrap(e))
				errCount++
				break
			}

			if size < limit { // ログを取ってるときに異常が発生してた。
				errCount++
				continue
			}

			// ファイルを回す。
			if e := rotateFile(path, num); e != nil {
				fmt.Fprintln(os.Stderr, erro.Wrap(e))
				errCount++
				break
			}

			errCount = 0
		}

		// 異常でループを抜けたら、デッドロック防止処理だけする。
		for {
			ent := <-queue
			if ent == FlushTrigger() {
				done <- false
			}
		}
	}, queueCapacity)
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
		if e := os.Rename(from, to); e != nil {
			return erro.Wrap(e)
		}
	}

	// 最新版 を .1 に。
	if e := os.Rename(path, path+".1"); e != nil {
		return erro.Wrap(e)
	}

	return nil
}
