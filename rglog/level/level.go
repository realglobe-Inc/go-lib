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

var lvToLabel []string = []string{
	OFF:   "OFF",
	ERR:   "ERR",
	WARN:  "WARN",
	INFO:  "INFO",
	DEBUG: "DEBUG",
	ALL:   "ALL",
}

var labelToLv map[string]Level

func init() {
	labelToLv = make(map[string]Level)
	for lv, label := range lvToLabel {
		labelToLv[label] = Level(lv)
	}
}

func (lv Level) String() string {
	val := int(lv)
	if val < 0 || len(lvToLabel) <= val {
		return "UNKNOWN"
	}
	return lvToLabel[val]
}

func ValueOf(label string) (Level, error) {
	lv, ok := labelToLv[label]
	if ok {
		return lv, nil
	} else {
		return 0, erro.New("level '", label, "' is not exist.")
	}
}
