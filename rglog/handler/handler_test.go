package handler

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"strconv"
	"testing"
	"time"
)

func testHandler(t *testing.T, hndl Handler) {
	// ただ走らせるだけ、確認できない。
	// panic で止まらないことの確認くらいにはなる。
	for _, lv := range level.Values() {
		hndl.SetLevel(lv)
		hndl.Output(&record{time.Now(), lv, "test", 0, lv.String()})
	}

	hndl.Flush()
	hndl.Close()
}

func benchmarkHandler(b *testing.B, hndl Handler) {
	defer hndl.Close()
	hndl.SetLevel(level.ALL)

	b.ResetTimer()
	date := time.Now()
	for i := 0; i < b.N; i++ {
		hndl.Output(&record{date.Add(time.Duration(i) * time.Nanosecond), level.INFO, "test", 0, strconv.Itoa(i)})
	}
}

type record struct {
	date time.Time
	lv   level.Level
	file string
	line int
	msg  string
}

func (rec *record) Date() time.Time {
	return rec.date
}
func (rec *record) Level() level.Level {
	return rec.lv
}
func (rec *record) File() string {
	return rec.file
}
func (rec *record) Line() int {
	return rec.line
}
func (rec *record) Message() string {
	return rec.msg
}
