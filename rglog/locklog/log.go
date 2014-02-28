package locklog

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"strings"
	"sync"
)

// 全部ロックするログ。

type logger struct {
	sync.Mutex
	name string

	level.Level
	hndls     map[handler.Handler]bool
	useParent bool
}

func (logg *logger) AddHandler(hndl handler.Handler) {
	logg.Lock()
	defer logg.Unlock()

	logg.hndls[hndl] = true
}

func (logg *logger) RemoveHandler(hndl handler.Handler) {
	logg.Lock()
	defer logg.Unlock()

	delete(logg.hndls, hndl)
}

func (logg *logger) SetLevel(lv level.Level) {
	logg.Lock()
	defer logg.Unlock()

	logg.Level = lv
}

func (logg *logger) SetUseParent(useParent bool) {
	logg.Lock()
	defer logg.Unlock()

	logg.useParent = useParent
}

func (logg *logger) Err(v ...interface{}) {
	logg.logging(level.ERR, v...)
}

func (logg *logger) Warn(v ...interface{}) {
	logg.logging(level.WARN, v...)
}

func (logg *logger) Info(v ...interface{}) {
	logg.logging(level.INFO, v...)
}

func (logg *logger) Debug(v ...interface{}) {
	logg.logging(level.DEBUG, v...)
}

func (logg *logger) logging(lv level.Level, v ...interface{}) {
	cur := logg

	cur.Lock()
	for {

		if lv <= cur.Level {
			for hndl, _ := range cur.hndls {
				hndl.Output(2, lv, v...)
			}
		}

		if !cur.useParent {
			cur.Unlock()
			return
		}

		lock.Lock() // ロック結合。
		cur.Unlock()

		// 結合する必要も無さそうだけど、念のため。
		// 結合しないなら getParent の中で lock.Lock() と defer lock.Unlock() すれば良い。
		// デッドロックを防ぐため、結合する順番は葉から根の方向のみ。

		cur = getParent(cur.name)

		if cur == nil {
			lock.Unlock()
			return
		}

		cur.Lock() // ロック結合。
		lock.Unlock()
	}
}

func (logg *logger) flush() {
	logg.Lock()
	defer logg.Unlock()

	for hndl, _ := range logg.hndls {
		hndl.Flush()
	}
}

// マップで仮想的に木構造を扱う。どうせ深さは 10 もいかない。

var lock sync.Mutex

var loggers map[string]*logger

const sep = "/"

// ロックは外で。
func getParent(name string) *logger {
	for curName := name; ; {
		pos := strings.LastIndex(curName, sep)
		if pos < 0 {
			// cur == github.com とか。
			return loggers[""]
		}

		parentName := curName[:pos]
		parent := loggers[parentName]
		if parent != nil {
			return parent
		}

		curName = parentName
	}
}

func init() {
	loggers = make(map[string]*logger)
}

func Logger(name string) *logger {
	lock.Lock()
	defer lock.Unlock()

	logg := loggers[name]
	if logg == nil {
		logg = &logger{name: name, useParent: true, hndls: make(map[handler.Handler]bool)}
		loggers[name] = logg
	}

	return logg
}

func Flush() {
	// デッドロックしないようにマップをさらってるときに logger 自体の処理はしない。
	lock.Lock()

	loggs := []*logger{}
	for _, logg := range loggers {
		loggs = append(loggs, logg)
	}

	lock.Unlock()

	for _, logg := range loggs {
		logg.flush()
	}
}
