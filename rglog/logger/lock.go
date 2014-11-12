package logger

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"strings"
	"sync"
)

// 全部ロックするログ。

type lockLogger struct {
	lock sync.Mutex
	name string

	lv        level.Level
	hndls     map[string]handler.Handler
	useParent bool

	mgr *lockManager // この lockLogger を作成した lockManager。
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

func (log *lockLogger) SetUseParent(useParent bool) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.useParent = useParent
}

func (log *lockLogger) Err(v ...interface{}) {
	log.logging(level.ERR, v...)
}

func (log *lockLogger) Warn(v ...interface{}) {
	log.logging(level.WARN, v...)
}

func (log *lockLogger) Info(v ...interface{}) {
	log.logging(level.INFO, v...)
}

func (log *lockLogger) Debug(v ...interface{}) {
	log.logging(level.DEBUG, v...)
}

func (log *lockLogger) logging(lv level.Level, v ...interface{}) {
	cur := log

	cur.lock.Lock()
	for {

		if lv <= cur.lv {
			for _, hndl := range cur.hndls {
				hndl.Output(2, lv, v...)
			}
		}

		if !cur.useParent {
			cur.lock.Unlock()
			return
		}

		cur.mgr.lock.Lock() // ロック結合。
		cur.lock.Unlock()

		// 結合する必要も無さそうだけど、念のため。
		// 結合しないなら getParent の中で lock.Lock() と defer lock.Unlock() すれば良い。
		// デッドロックを防ぐため、結合する順番は葉から根の方向のみ。

		newCur := cur.mgr.getParent(cur.name)
		if newCur == nil {
			cur.mgr.lock.Unlock()
			return
		}
		cur = newCur

		cur.lock.Lock() // ロック結合。
		cur.mgr.lock.Unlock()
	}
}

func (log *lockLogger) flush() {
	log.lock.Lock()
	defer log.lock.Unlock()

	for _, hndl := range log.hndls {
		hndl.Flush()
	}
}

type lockManager struct {
	lock sync.Mutex

	// マップで仮想的に木構造を扱う。どうせ深さは 10 もいかない。
	loggers map[string]*lockLogger
}

func NewLockManager() *lockManager {
	return &lockManager{loggers: map[string]*lockLogger{}}
}

// ロックは外で。
func (mgr *lockManager) getParent(name string) *lockLogger {
	const sep = "/"
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

func (mgr *lockManager) Logger(name string) Logger {
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

func (mgr *lockManager) Flush() {
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
