package level

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

type Level int

const (
	// 表示レベル無指定なら OFF (0) になるので、何も出力しない。
	OFF Level = iota

	ERR
	WARN
	INFO
	DEBUG

	ALL
)

var labels []string = []string{
	OFF:   "OFF",
	ERR:   "ERR",
	WARN:  "WARN",
	INFO:  "INFO",
	DEBUG: "DEBUG",
	ALL:   "ALL",
}

func (lv Level) String() string {
	if int(lv) < 0 || len(labels) <= int(lv) {
		return "UNKNOWN"
	}
	return labels[lv]
}

func ValueOf(str string) (Level, error) {
	for lv, label := range labels {
		if str == label {
			return Level(lv), nil
		}
	}
	return 0, erro.New("level '", str, "' is not exist.")
}
