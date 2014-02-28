package handler

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"time"
)

type Formatter interface {
	Format(date time.Time, file string, line int, lv level.Level, v ...interface{}) []byte
}

type simpleFormatter struct {
}

func (formatter simpleFormatter) Format(date time.Time, file string, line int, lv level.Level, v ...interface{}) []byte {
	year, month, day := date.Date()
	hour, min, sec := date.Clock()
	microSec := date.Nanosecond() / 1000

	// Level は固定幅で左寄せ。見た目はいいけど cut では後ろが分けられない。
	msg := fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d.%06d %-5v %s:%d ", year, int(month), day, hour, min, sec, microSec, lv, file, line) +
		fmt.Sprint(v...) +
		"\n"

	return []byte(msg)
}
