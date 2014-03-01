package handler

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"os"
	"strings"
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

	// コンパイル環境の ${GOPATH}/src/ を取り除く。
	// os.Getenv("GOPATH") して得られるのは実行環境の値なので、不完全な方法を取る。
	srcDir := string(os.PathSeparator) + "src" + string(os.PathSeparator)
	pos := strings.Index(file, srcDir)
	if pos >= 0 {
		file = file[pos+len(srcDir):]
	}

	// Level は固定幅で左寄せ。見た目はいいけど cut では後ろが分けられない。
	msg := fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d.%06d %-5v %s:%d ", year, int(month), day, hour, min, sec, microSec, lv, file, line) +
		fmt.Sprint(v...) +
		"\n"

	return []byte(msg)
}
