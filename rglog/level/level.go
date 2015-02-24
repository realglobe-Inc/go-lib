package level

import (
	"github.com/realglobe-Inc/go-lib/erro"
)

type Level int

const (
	// 表示レベル無指定 0 なら OFF になるので、何も出力しない。
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

// lv が lv2 より重要なときのみ true。
func (lv Level) Higher(lv2 Level) bool {
	return lv < lv2
}

// lv が lv2 より重要でないときのみ true。
func (lv Level) Lower(lv2 Level) bool {
	return lv > lv2
}

func (lv Level) String() string {
	val := int(lv)
	if val < 0 || len(lvToLabel) <= val {
		return "UNKNOWN"
	}
	return lvToLabel[val]
}

// 文字列から値に。
func ValueOf(label string) (Level, error) {
	lv, ok := labelToLv[label]
	if ok {
		return lv, nil
	} else {
		return 0, erro.New("level " + label + " is not exist.")
	}
}

// 重要度降順で列挙する。
func Values() []Level {
	values := []Level{}
	for i := OFF; i <= ALL; i++ {
		values = append(values, i)
	}
	return values
}

// flag.Var で使う。
type levelVar struct {
	*Level
}

func (v levelVar) Set(s string) error {
	var err error
	*v.Level, err = ValueOf(s)
	if err != nil {
		return erro.Wrap(err)
	}
	return nil
}

// flags.Var(Var(&param, level.INFO), "level", "Log level.") の形で使う。
func Var(lv *Level, defaultLv Level) levelVar {
	*lv = defaultLv
	return levelVar{lv}
}
