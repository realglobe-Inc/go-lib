package handler

import (
	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
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

func TestRotateHandlerLevel(t *testing.T) {
	path, err := testFilePath()
	if err != nil {
		t.Fatal(err)
	}
	hndl := NewRotateHandler(path, 1<<20, 10)
	defer os.Remove(path)
	for i := 1; i <= 10; i++ {
		defer os.Remove(path + "." + strconv.Itoa(i))
	}

	testHandlerLevel(t, hndl)
}

func TestRotateHandlerOutput(t *testing.T) {
	path, err := testFilePath()
	if err != nil {
		t.Fatal(err)
	}
	hndl := NewRotateHandler(path, 1<<20, 10)
	defer os.Remove(path)
	for i := 1; i <= 10; i++ {
		defer os.Remove(path + "." + strconv.Itoa(i))
	}

	testHandlerOutput(t, hndl)
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
	defer hndl.Close()

	for i := 0; i < 2*num; i++ {
		hndl.Output(&record{date: time.Now(), lv: level.INFO, msg: "test"})
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
	defer hndl.Close()

	n := 1000
	for i := 0; i < n; i++ {
		hndl.Output(&record{date: time.Now(), lv: level.INFO, msg: "test"})
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
