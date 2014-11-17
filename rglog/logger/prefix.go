package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ログにコンパイル環境の恥ずかしいパスとかを載せないように。

// ${GOPATH}/src/ の部分。
var uselessPrefix string

func init() {
	// このファイルの名前から ${GOPATH}/src/ を逆算する。

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return
	}

	suffix := filepath.Join("github.com", "realglobe-Inc", "go-lib-rg", "rglog", "logger", "prefix.go")
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
