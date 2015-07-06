// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog/level"
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
			t.Fatal(bak, err)
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
		t.Fatal(n, len(lines))
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
