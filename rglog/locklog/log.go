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

func (log *logger) AddHandler(hndl handler.Handler) {
	log.Lock()
	defer log.Unlock()

	log.hndls[hndl] = true
}

func (log *logger) RemoveHandler(hndl handler.Handler) {
	log.Lock()
	defer log.Unlock()

	delete(log.hndls, hndl)
}

func (log *logger) SetLevel(lv level.Level) {
	log.Lock()
	defer log.Unlock()

	log.Level = lv
}

func (log *logger) SetUseParent(useParent bool) {
	log.Lock()
	defer log.Unlock()

	log.useParent = useParent
}

func (log *logger) Err(v ...interface{}) {
	log.logging(level.ERR, v...)
}

func (log *logger) Warn(v ...interface{}) {
	log.logging(level.WARN, v...)
}

func (log *logger) Info(v ...interface{}) {
	log.logging(level.INFO, v...)
}

func (log *logger) Debug(v ...interface{}) {
	log.logging(level.DEBUG, v...)
}

func (log *logger) logging(lv level.Level, v ...interface{}) {
	cur := log

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

func (log *logger) flush() {
	log.Lock()
	defer log.Unlock()

	for hndl, _ := range log.hndls {
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

	log := loggers[name]
	if log == nil {
		log = &logger{name: name, useParent: true, hndls: make(map[handler.Handler]bool)}
		loggers[name] = log
	}

	return log
}

func Flush() {
	// デッドロックしないようにマップをさらってるときに logger 自体の処理はしない。
	lock.Lock()

	logs := []*logger{}
	for _, log := range loggers {
		logs = append(logs, log)
	}

	lock.Unlock()

	for _, log := range logs {
		log.flush()
	}
}
