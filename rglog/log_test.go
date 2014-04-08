package rglog

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	rootLabel := "github.com/realglobe-Inc/go-lib-rg"

	loop := 100000
	n := 100

	rootLog := GetLogger(rootLabel)
	rootLog.SetLevel(level.ALL)
	rootLog.SetUseParent(false)

	path := filepath.Join(os.TempDir(), "log_test.go.log")
	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}
	defer os.Remove(path)

	hndl, err := handler.NewRotateHandler(path, 1<<30, 10)
	if err != nil {
		t.Fatal(err)
	}
	hndl.SetLevel(level.DEBUG)
	rootLog.AddHandler(hndl)

	for i := 0; i < loop; i++ {
		GetLogger(rootLabel + "/" + strconv.Itoa(i%n)).Info(i)
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
	rootLog.SetLevel(level.ALL)
	rootLog.SetUseParent(false)

	path := filepath.Join(os.TempDir(), "log_test.go.log")
	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}
	defer os.Remove(path)

	hndl, err := handler.NewRotateHandler(path, 1<<30, 10)
	if err != nil {
		t.Fatal(err)
	}
	hndl.SetLevel(level.DEBUG)
	rootLog.AddHandler(hndl)

	c := make(chan bool)
	timeout := time.After(time.Duration(int64(n*loop*100) * int64(time.Microsecond)))

	for i := 0; i < n; i++ {
		id := i
		go func() {
			for j := 0; j < loop; j++ {
				GetLogger(rootLabel+"/"+strconv.Itoa(id)).Info(id, j)
			}

			c <- true
		}()
	}

	for i := 0; i < n; i++ {
		select {
		case <-c:
		case <-timeout:
			t.Fatal("Dead lock?")
		}
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
