package handler

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"os"
	"path/filepath"
	"runtime"
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

	if strings.HasPrefix(file, uselessPrefix) { // 違う環境でコンパイルした後、リンクすることは可能だと思うので。
		file = file[len(uselessPrefix):]
	} else {
		// /src/ の前までを GOPATH とみなす。
		// GOPATH 自体に /src/ が含まれていると、そこまでしか除去できない。
		srcDir := string(os.PathSeparator) + "src" + string(os.PathSeparator)
		pos := strings.Index(file, srcDir)
		if pos >= 0 {
			file = file[pos+len(srcDir):]
		}
	}

	// Level は固定幅で左寄せ。見た目はいいけど cut では後ろが分けられない。
	msg := fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d.%06d %-5v %s:%d ",
		year, int(month), day, hour, min, sec, microSec, lv, file, line) +
		fmt.Sprint(v...) + "\n"

	return []byte(msg)
}

// ${GOPATH}/src/ の部分。
var uselessPrefix string

func init() {
	// このファイルの名前から ${GOPATH}/src/ を逆算する。

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return
	}

	suffix := filepath.Join("github.com", "realglobe-Inc", "go-lib-rg", "rglog", "handler", "formatter.go")
	if !strings.HasSuffix(file, suffix) {
		return
	}

	uselessPrefix = file[:len(file)-len(suffix)]
}
