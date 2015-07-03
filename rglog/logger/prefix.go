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

package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ログにコンパイル環境の恥ずかしいパスとかを載せないように。

// ${GOPATH}/src/ の部分。
var uselessPref string

func init() {
	// このファイルの名前から ${GOPATH}/src/ を逆算する。

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return
	}

	suffix := filepath.Join("github.com", "realglobe-Inc", "go-lib", "rglog", "logger", "prefix.go")
	if !strings.HasSuffix(file, suffix) {
		return
	}

	uselessPref = file[:len(file)-len(suffix)]
}

// ファイル名から ${GOPATH}/src/ の部分を除く。
func trimPrefix(file string) string {
	if uselessPref != "" && strings.HasPrefix(file, uselessPref) { // 違う環境でコンパイルした後、リンクすることは可能だと思うので。
		return file[len(uselessPref):]
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
