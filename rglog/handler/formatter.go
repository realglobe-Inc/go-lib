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

type simpleFormatter struct{}

var SimpleFormatter = &simpleFormatter{}

func (formatter simpleFormatter) Format(date time.Time, file string, line int, lv level.Level, v ...interface{}) []byte {
	year, month, day := date.Date()
	hour, min, sec := date.Clock()
	microSec := date.Nanosecond() / 1000

	file = trimPrefix(file)

	// Level は最短長固定幅。
	msg := fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d.%06d %.3v %s:%d %s\n",
		year, int(month), day, hour, min, sec, microSec, lv, file, line, fmt.Sprint(v...))

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

// ファイル名から ${GOPATH}/src/ の部分を除く。
func trimPrefix(file string) string {
	if uselessPrefix != "" && strings.HasPrefix(file, uselessPrefix) { // 違う環境でコンパイルした後、リンクすることは可能だと思うので。
		return file[len(uselessPrefix):]
	} else {
		// /src/ の前までを GOPATH とみなす。
		// GOPATH 自体に /src/ が含まれていると、そこまでしか除去できない。
		srcDir := string(os.PathSeparator) + "src" + string(os.PathSeparator)
		pos := strings.Index(file, srcDir)
		if pos >= 0 {
			return file[pos+len(srcDir):]
		}
	}
	return file
}

type levelOnlyFormatter struct{}

var LevelOnlyFormatter = &levelOnlyFormatter{}

func (formatter levelOnlyFormatter) Format(date time.Time, file string, line int, lv level.Level, v ...interface{}) []byte {
	msg := fmt.Sprintf("[%.3v] %s\n", lv, fmt.Sprint(v...))
	return []byte(msg)
}
