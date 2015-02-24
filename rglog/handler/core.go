package handler

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"
)

// ログの書き込みを別ゴルーチンで実行できるようにするために分離。
// synchronizedCoreHandler でラップして使うので、この部分をスレッドセーフに実装する必要はない。
type coreHandler interface {
	output(rec Record)
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
	rec Record
}

type synchronizedFlushRequest struct {
	ackCh chan<- struct{}
}

type synchronizedCloseRequest struct {
	ackCh chan<- struct{}
}

func (core *synchronizedCoreHandler) output(rec Record) {
	core.reqCh <- &synchronizedOutputRequest{rec}
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
		base.output(r.rec)
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
	// synchronizedCoreHandler でラップしておけば問題無い。
	base coreHandler
}

func wrapCoreHandler(core coreHandler) Handler {
	return &coreWrapper{lv: level.ALL, base: core}
}

func (hndl *coreWrapper) Level() level.Level {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	return hndl.lv
}

func (hndl *coreWrapper) SetLevel(lv level.Level) {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	hndl.lv = lv
}

func (hndl *coreWrapper) Output(rec Record) {
	hndl.lock.Lock()
	if rec.Level().Lower(hndl.lv) {
		hndl.lock.Unlock()
		return
	}
	hndl.lock.Unlock()

	hndl.base.output(rec)
}

func (hndl *coreWrapper) Flush() {
	hndl.base.flush()
}

func (hndl *coreWrapper) Close() {
	hndl.base.close()
}
