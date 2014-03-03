package rglog

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	rootLabel := "github.com/realglobe-Inc/go-lib-rg"

	loop := 100000
	n := 100

	rootLog := GetLogger(rootLabel)
	rootLog.SetLevel(level.DEBUG)
	rootLog.SetUseParent(false)

	hndl := handler.NewConsoleHandler()
	hndl.SetLevel(level.INFO)
	rootLog.AddHandler(hndl)

	path := filepath.Join(os.TempDir(), "log_test.go.log")
	if e := os.Remove(path); e != nil {
		if !os.IsNotExist(e) {
			t.Fatal(e)
		}
	}
	hndl, err := handler.NewFileHandler(path)
	if err != nil {
		t.Fatal(err)
	}
	hndl.SetLevel(level.DEBUG)
	rootLog.AddHandler(hndl)

	start := time.Now()
	for i := 0; i < loop; i++ {
		GetLogger(rootLabel + "/" + strconv.Itoa(i%n)).Info(i)
	}
	end := time.Now()

	// 遅過ぎ検知。
	// 1 回 100 マイクロ秒も掛かってるのは遅い。
	limit := start.Add(time.Duration(int64(loop*100) * int64(time.Microsecond)))
	if end.After(limit) {
		t.Error("Too slow ", end.Sub(start))
	} else {
		//t.Error("Not too slow", end.Sub(start))
	}

	Flush()

	// ファイルに書き込めているかどうか検査。
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(buff) > 0 && buff[len(buff)-1] == '\n' { // 最後の空行は抜かしとく。
		buff = buff[:len(buff)-1]
	}

	lines := strings.Split(string(buff), "\n")
	if len(lines) != loop {
		t.Error(len(lines), loop)
	}

}

func TestConcurrent(t *testing.T) {
	rootLabel := "github.com/realglobe-Inc/go-lib-rg"

	n := 100
	loop := 1000

	rootLog := GetLogger(rootLabel)
	rootLog.SetLevel(level.DEBUG)
	rootLog.SetUseParent(false)

	hndl := handler.NewConsoleHandler()
	hndl.SetLevel(level.INFO)
	rootLog.AddHandler(hndl)

	path := filepath.Join(os.TempDir(), "log_test.go.log")
	if e := os.Remove(path); e != nil {
		if !os.IsNotExist(e) {
			t.Fatal(e)
		}
	}
	hndl, err := handler.NewFileHandler(path)
	if err != nil {
		t.Fatal(err)
	}
	hndl.SetLevel(level.DEBUG)
	rootLog.AddHandler(hndl)

	var lock sync.Mutex

	var zero time.Time
	start := time.Now()
	end := zero

	c := make(chan bool)

	for i := 0; i < n; i++ {
		id := i
		go func() {
			for j := 0; j < loop; j++ {
				GetLogger(rootLabel+"/"+strconv.Itoa(id)).Info(id, j)
			}

			c <- true
		}()
	}

	go func() {
		for i := 0; i < n; i++ {
			<-c
		}

		lock.Lock()
		end = time.Now()
		lock.Unlock()
	}()

	// 遅過ぎ検知。
	// 1 回 100 マイクロ秒も掛かってるのは遅い。
	limit := start.Add(time.Duration(int64(n*loop*100) * int64(time.Microsecond)))
	for time.Now().Before(limit) {
		lock.Lock()
		curEnd := end
		lock.Unlock()

		if curEnd != zero {
			break
		}

		time.Sleep(time.Millisecond)
	}

	if end == zero {
		t.Fatal("Too slow ", time.Now().Sub(start))
	} else {
		//t.Error("Not too slow", end.Sub(start))
	}

	Flush()

	// ファイルに書き込めているかどうか検査。
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(buff) > 0 && buff[len(buff)-1] == '\n' { // 最後の空行は抜かしとく。
		buff = buff[:len(buff)-1]
	}

	lines := strings.Split(string(buff), "\n")
	if len(lines) != n*loop {
		t.Error(len(lines), n*loop)
	}

}
