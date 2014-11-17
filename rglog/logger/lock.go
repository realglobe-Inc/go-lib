package logger

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 全部ロックするログ。

type lockLogger struct {
	lock sync.Mutex
	name string

	lv        level.Level
	hndls     map[string]handler.Handler
	useParent bool

	mgr *lockLoggerManager // この lockLogger を作成した lockLoggerManager。
}

func (log *lockLogger) Handler(key string) handler.Handler {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.hndls[key]
}

func (log *lockLogger) AddHandler(key string, hndl handler.Handler) handler.Handler {
	log.lock.Lock()
	defer log.lock.Unlock()

	old := log.hndls[key]
	log.hndls[key] = hndl
	return old
}

func (log *lockLogger) RemoveHandler(key string) handler.Handler {
	log.lock.Lock()
	defer log.lock.Unlock()

	old := log.hndls[key]
	delete(log.hndls, key)
	return old
}

func (log *lockLogger) Level() level.Level {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.lv
}

func (log *lockLogger) SetLevel(lv level.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.lv = lv
}

func (log *lockLogger) UseParent() bool {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.useParent
}

func (log *lockLogger) SetUseParent(useParent bool) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.useParent = useParent
}

func (log *lockLogger) IsLoggable(lv level.Level) bool {
	cur := log

	for {
		cur.lock.Lock()
		curLv := cur.lv
		hndlNum := len(cur.hndls)
		useParent := cur.useParent
		cur.lock.Unlock()

		if lv <= curLv && hndlNum > 0 {
			return true
		}

		if !useParent {
			return false
		}

		newCur := cur.mgr.getParent(cur.name)
		if newCur == nil {
			return false
		}
		cur = newCur
	}
}

func (log *lockLogger) logging(rec *record) {
	cur := log

	for {
		hndls := []handler.Handler{}
		cur.lock.Lock()
		lv := cur.lv
		for _, hndl := range cur.hndls {
			hndls = append(hndls, hndl)
		}
		useParent := cur.useParent
		cur.lock.Unlock()

		if rec.Level() <= lv && len(hndls) > 0 {
			if rec.file == "" {
				rec.date = time.Now()
				if _, file, line, ok := runtime.Caller(2); ok {
					rec.file = trimPrefix(file)
					rec.line = line
				} else {
					rec.file = "???"
					rec.line = 0
				}
				rec.msg = fmt.Sprint(rec.rawMsg...)
			}

			for _, hndl := range hndls {
				hndl.Output(rec)
			}
		}

		if !useParent {
			return
		}

		// ロック結合した方が良さそうだけど、たぶん大丈夫だろう。
		newCur := cur.mgr.getParent(cur.name)
		if newCur == nil {
			return
		}
		cur = newCur
	}
}

func (log *lockLogger) Log(lv level.Level, v ...interface{}) {
	log.logging(&record{lv: lv, rawMsg: v})
}

func (log *lockLogger) Err(v ...interface{}) {
	log.logging(&record{lv: level.ERR, rawMsg: v})
}

func (log *lockLogger) Warn(v ...interface{}) {
	log.logging(&record{lv: level.WARN, rawMsg: v})
}

func (log *lockLogger) Info(v ...interface{}) {
	log.logging(&record{lv: level.INFO, rawMsg: v})
}

func (log *lockLogger) Debug(v ...interface{}) {
	log.logging(&record{lv: level.DEBUG, rawMsg: v})
}

func (log *lockLogger) flush() {
	log.lock.Lock()
	defer log.lock.Unlock()

	for _, hndl := range log.hndls {
		hndl.Flush()
	}
}

type lockLoggerManager struct {
	lock sync.Mutex

	// マップで仮想的に木構造を扱う。どうせ深さは 10 もいかない。
	loggers map[string]*lockLogger
}

func NewLockLoggerManager() *lockLoggerManager {
	return &lockLoggerManager{loggers: map[string]*lockLogger{}}
}

// ロックは外で。
func (mgr *lockLoggerManager) getParent(name string) *lockLogger {
	const sep = "/"

	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	for curName := name; ; {
		pos := strings.LastIndex(curName, sep)
		if pos < 0 {
			// cur == github.com とか。
			return mgr.loggers[""]
		}

		parentName := curName[:pos]
		parent := mgr.loggers[parentName]
		if parent != nil {
			return parent
		}

		curName = parentName
	}
}

func (mgr *lockLoggerManager) Logger(name string) Logger {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	log := mgr.loggers[name]
	if log == nil {
		log = &lockLogger{
			name:      name,
			useParent: true,
			hndls:     make(map[string]handler.Handler),
			mgr:       mgr,
		}
		mgr.loggers[name] = log
	}

	return log
}

func (mgr *lockLoggerManager) Flush() {
	// デッドロックしないようにマップをさらってるときに lockLogger 自体の処理はしない。
	mgr.lock.Lock()

	logs := []*lockLogger{}
	for _, log := range mgr.loggers {
		logs = append(logs, log)
	}

	mgr.lock.Unlock()

	for _, log := range logs {
		log.flush()
	}
}

type record struct {
	date   time.Time
	lv     level.Level
	file   string
	line   int
	msg    string
	rawMsg []interface{}
}

func (rec *record) Date() time.Time {
	return rec.date
}
func (rec *record) Level() level.Level {
	return rec.lv
}
func (rec *record) File() string {
	return rec.file
}
func (rec *record) Line() int {
	return rec.line
}
func (rec *record) Message() string {
	return rec.msg
}
