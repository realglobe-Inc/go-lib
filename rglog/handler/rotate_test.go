package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

func testFilePath() (path string, err error) {
	file, err := ioutil.TempFile("", "test_")
	if err != nil {
		return "", erro.Wrap(err)
	}
	if err := file.Close(); err != nil {
		return "", erro.Wrap(err)
	}
	if err := os.Remove(file.Name()); err != nil {
		return "", erro.Wrap(err)
	}
	return file.Name(), nil
}

// ローテートしてみる。
func TestRotateHandlerRotation(t *testing.T) {

	path, err := testFilePath()
	if err != nil {
		t.Fatal(err)
	}
	num := 100
	hndl := NewRotateHandler(path, 0, num)
	defer os.Remove(path)
	for i := 1; i <= num; i++ {
		defer os.Remove(path + "." + strconv.Itoa(i))
	}

	hndl.SetLevel(level.ALL)
	for i := 0; i < 2*num; i++ {
		hndl.Output(0, level.INFO, "test")
	}
	hndl.Flush()

	for i := 1; i <= num; i++ {
		bak := path + "." + strconv.Itoa(i)
		if _, err := os.Stat(bak); err != nil {
			t.Error(bak, err)
		}
	}
}

// 記録する。
func TestRotateHandlerLogging(t *testing.T) {

	path, err := testFilePath()
	if err != nil {
		t.Fatal(err)
	}
	hndl := NewRotateHandler(path, 1<<20, 0)
	defer os.Remove(path)

	n := 1000
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

func BenchmarkRotateHandler(b *testing.B) {
	path, err := testFilePath()
	if err != nil {
		b.Fatal(err)
	}
	hndl := NewRotateHandler(path, 1<<20, 10)
	defer os.Remove(path)
	for i := 1; i <= 10; i++ {
		defer os.Remove(path + "." + strconv.Itoa(i))
	}

	benchmarkHandler(b, hndl)
}
