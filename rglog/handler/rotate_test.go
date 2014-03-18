package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

// ローテートしてみる。
func TestRotateHandlerRotation(t *testing.T) {

	file, err := ioutil.TempFile("", "go_rotate_test")
	if err != nil {
		t.Fatal(err)
	}
	if e := file.Close(); e != nil {
		t.Fatal(e)
	}
	if e := os.Remove(file.Name()); e != nil {
		t.Fatal(e)
	}
	defer os.Remove(file.Name())

	path := file.Name()
	n := 100

	hndl := NewRotateHandler(path, 0, n)
	hndl.SetLevel(level.ALL)

	for i := 0; i < n; i++ {
		hndl.Output(0, level.INFO, "test")
	}

	hndl.Flush()

	for i := 1; i <= n; i++ {
		bak := path + "." + strconv.Itoa(i)
		if _, e := os.Stat(bak); e != nil {
			t.Error(bak, e)
		}

		os.Remove(bak)
	}

}

// 記録する。
func TestRotateHandlerLogging(t *testing.T) {

	file, err := ioutil.TempFile("", "go_rotate_test")
	if err != nil {
		t.Fatal(err)
	}
	if e := file.Close(); e != nil {
		t.Fatal(e)
	}
	if e := os.Remove(file.Name()); e != nil {
		t.Fatal(e)
	}
	defer os.Remove(file.Name())

	path := file.Name()
	n := 100

	hndl := NewRotateHandler(path, 1<<20, 10)
	hndl.SetLevel(level.ALL)

	for i := 0; i < n; i++ {
		hndl.Output(0, level.INFO, "test")
	}

	hndl.Flush()

	buff, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(buff) > 0 && buff[len(buff)-1] == '\n' {
		buff = buff[:len(buff)-1]
	}

	lines := strings.Split(string(buff), "\n")
	if len(lines) != n {
		t.Error(n, len(lines))
	}
}
