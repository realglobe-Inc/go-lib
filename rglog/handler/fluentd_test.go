package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"net"
	"testing"
	"time"
)

// テストしたかったら fluentd サーバーを立ててから。
var fluentdAddr = "localhost:24224"

func init() {
	if fluentdAddr != "" {
		// 実際にサーバーが立っているかどうか調べる。
		// 立ってなかったらテストはスキップ。
		conn, err := net.Dial("tcp", fluentdAddr)
		if err != nil {
			fluentdAddr = ""
		} else {
			conn.Close()
		}
	}
}

func TestFluentdHandlerLevel(t *testing.T) {
	if fluentdAddr == "" {
		t.SkipNow()
	}

	testHandlerLevel(t, NewFluentdHandler(fluentdAddr, "rglog.test"))
}

func TestFluentdHandlerOutput(t *testing.T) {
	if fluentdAddr == "" {
		t.SkipNow()
	}

	testHandlerOutput(t, NewFluentdHandler(fluentdAddr, "rglog.test"))
}

// 色んな長さのメッセージを送る。
// MessagePack 部分のテスト。
func TestFluentdHandlerMessageLength(t *testing.T) {
	if fluentdAddr == "" {
		t.SkipNow()
	}

	hndl := NewFluentdHandler(fluentdAddr, "rglog.test")
	defer hndl.Close()

	for i := 0; i < 18; i++ {
		msg := ""
		for j := 0; j < (1 << uint(i)); j++ {
			msg += "a"
		}
		hndl.Output(&record{date: time.Now(), lv: level.INFO, msg: msg})
	}
}

// 複数接続で。
func TestManyFluentdHandler(t *testing.T) {
	if fluentdAddr == "" {
		t.SkipNow()
	}

	n := 20
	loop := 100

	hndls := []Handler{}
	for i := 0; i < n; i++ {
		hndl := NewFluentdHandler(fluentdAddr, "rglog.test")
		defer hndl.Close()

		hndls = append(hndls, hndl)
	}

	for i := 0; i < loop; i++ {
		for j := 0; j < len(hndls); j++ {
			hndls[j].Output(&record{date: time.Now(), lv: level.INFO, msg: fmt.Sprint("a ", j, i)})
		}
	}
}

// サーバー側からの切断に耐えるか。
func TestFluentdHandlerConnectionCut(t *testing.T) {
	// A から B にログ 1 を流す。
	// B はログ 1 を受け取ったら切断する。
	// 同期。
	// A から B にログ 2 を流す。
	// A から B にログ 3 を流す。
	// B はログ 2 とログ3 を受け取る。
	// 同期。

	// サーバー側で閉じられた接続に対する 1 回目の書き込みではエラーが出ないので、
	// 切断後に 2 回書き込む必要がある。

	syncCh := make(chan error)

	// A.
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer lis.Close()
	go func() {
		conn, err := lis.Accept()
		if err != nil {
			syncCh <- err
			return
		}

		buff := make([]byte, 8192)
		if n, err := conn.Read(buff); err != nil {
			syncCh <- err
			return
		} else if bytes.Index(buff[:n], []byte("log message 1")) < 0 {
			syncCh <- errors.New("no log message 1")
		}

		// ログ 1 を受け取った。

		if err := conn.Close(); err != nil {
			syncCh <- err
			return
		}
		conn = nil
		time.Sleep(10 * time.Millisecond)

		// 切断した。

		syncCh <- nil

		// 同期。

		conn, err = lis.Accept()
		if err != nil {
			syncCh <- err
			return
		}

		if n, err := conn.Read(buff); err != nil {
			syncCh <- err
			return
		} else if bytes.Index(buff[:n], []byte("log message 2")) < 0 {
			syncCh <- errors.New("no log message 2")
		} else if bytes.Index(buff[:n], []byte("log message 3")) < 0 {
			syncCh <- errors.New("no log message 3")
		}

		// ログ 2 を受け取った。

		syncCh <- nil
	}()

	// B.
	// A の起動待ち。
	time.Sleep(10 * time.Millisecond)

	hndl := NewFluentdHandler(lis.Addr().String(), "rglog.test")
	defer hndl.Close()

	hndl.Output(&record{date: time.Now(), lv: level.INFO, msg: "log message 1"})
	hndl.Flush()

	// ログ 1 を流した。

	if err := <-syncCh; err != nil {
		t.Fatal(err)
	}

	// 同期。

	hndl.Output(&record{date: time.Now(), lv: level.INFO, msg: "log message 2"})
	hndl.Flush()

	hndl.Output(&record{date: time.Now(), lv: level.INFO, msg: "log message 3"})
	hndl.Flush()

	if err := <-syncCh; err != nil {
		t.Fatal(err)
	}
}

func BenchmarkFluentdHandler(b *testing.B) {
	if fluentdAddr == "" {
		b.SkipNow()
	}

	benchmarkHandler(b, NewFluentdHandler(fluentdAddr, "rglog.test"))
}
