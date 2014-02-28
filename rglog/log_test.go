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
	end := false

	start := time.Now()

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
		end = true
		lock.Unlock()
	}()

	// 遅過ぎ検知。
	// 1 回 0.1 ミリ秒も掛かってるのは遅い。
	limit := start.Add(time.Duration(n * loop * int(time.Millisecond) / 10))
	for time.Now().Before(limit) {
		lock.Lock()
		flag := end
		lock.Unlock()

		if flag {
			break
		}

		time.Sleep(time.Millisecond)
	}

	if !end {
		t.Fatal("Too slow")
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
