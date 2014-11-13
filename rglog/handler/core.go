package handler

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"
)

// ログの書き込みを別ゴルーチンで実行できるようにするために分離。
// 別ゴルーチンだとファイル名と行番号の取得ができないので、こんな切り分け。
// synchronizedCoreHandler でラップして使うので、この部分をスレッドセーフに実装する必要はない。
type coreHandler interface {
	output(file string, line int, lv level.Level, v ...interface{})
	flush()
	close()
}

// coreHandler をスレッドセーフにするラッパー。
// ついでに別ゴルーチンでの書き込みにもなる。
// output はノンブロッキング。
// flush, close はブロッキング。
type synchronizedCoreHandler struct {
	reqCh chan<- interface{}
}

// 書き出し待機させる最大数。
const chCap = 1000

// やることが無いときに flush する間隔。
const flushInterval = time.Minute

type synchronizedOutputRequest struct {
	file string
	line int
	level.Level
	v []interface{}
}

type synchronizedFlushRequest struct {
	ackCh chan<- struct{}
}

type synchronizedCloseRequest struct {
	ackCh chan<- struct{}
}

func (core *synchronizedCoreHandler) output(file string, line int, lv level.Level, v ...interface{}) {
	core.reqCh <- &synchronizedOutputRequest{file, line, lv, v}
}

func (core *synchronizedCoreHandler) flush() {
	ackCh := make(chan struct{}, 1)
	core.reqCh <- &synchronizedFlushRequest{ackCh}
	<-ackCh
}

func (core *synchronizedCoreHandler) close() {
	ackCh := make(chan struct{}, 1)
	core.reqCh <- &synchronizedCloseRequest{ackCh}
	<-ackCh
}

func newSynchronizedCoreHandler(base coreHandler) coreHandler {
	reqCh := make(chan interface{}, chCap)

	go func() {
		closed := false
		for !closed {
			func() { // パニックになったときも素知らぬ顔で次のリクエストを処理するために関数で括る。
				defer func() {
					if rcv := recover(); rcv != nil {
						buff := make([]byte, 8192)
						stackLen := runtime.Stack(buff, false)
						stack := string(buff[:stackLen])

						fmt.Fprintln(os.Stderr, rcv)
						fmt.Fprintln(os.Stderr, stack)
					}
				}()

				select {
				case req := <-reqCh:
					closed = handleSynchronizedRequest(base, req)
					return
				default:
				}

				// 毎回タイマー設定するのが気持ち悪いので、本当に待つときだけ。
				timer := time.NewTimer(flushInterval)
				defer timer.Stop()

				select {
				case req := <-reqCh:
					closed = handleSynchronizedRequest(base, req)
				case <-timer.C:
					base.flush()
				}
			}()
		}
	}()

	return &synchronizedCoreHandler{reqCh}
}

func handleSynchronizedRequest(base coreHandler, req interface{}) (closed bool) {
	switch r := req.(type) {
	case *synchronizedOutputRequest:
		base.output(r.file, r.line, r.Level, r.v...)
	case *synchronizedFlushRequest:
		defer func() { r.ackCh <- struct{}{} }()
		base.flush()
	case *synchronizedCloseRequest:
		defer func() { r.ackCh <- struct{}{} }()
		base.close()
		return true
	default:
		panic("unknown request " + reflect.TypeOf(req).Name())
	}
	return false
}

// coreHandler をラップして Handler にする。
// ファイル名と行番号を取得しつつスレッドセーフにするためにロックが必要。
type coreWrapper struct {
	lock sync.Mutex
	lv   level.Level

	// スレッドセーフに使いたいなら、こいつをスレッドセーフにしとく必要あり。
	base coreHandler
}

func wrapCoreHandler(core coreHandler) Handler {
	return &coreWrapper{base: core}
}

func (hndl *coreWrapper) SetLevel(lv level.Level) {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	hndl.lv = lv
}

func (hndl *coreWrapper) Output(depth int, lv level.Level, v ...interface{}) {
	hndl.lock.Lock()
	if lv > hndl.lv {
		hndl.lock.Unlock()
		return
	}
	hndl.lock.Unlock()

	_, file, line, ok := runtime.Caller(depth + 1)
	if ok {
		file = trimPrefix(file)
	} else {
		file = "???"
		line = 0
	}

	hndl.base.output(file, line, lv, v...)
}

func (hndl *coreWrapper) Flush() {
	hndl.base.flush()
}

func (hndl *coreWrapper) Close() {
	hndl.base.close()
}
