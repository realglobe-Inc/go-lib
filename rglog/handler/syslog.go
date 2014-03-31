package handler

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"log/syslog"
	"runtime"
	"sync"
)

type SyslogHandler struct {
	// 隠蔽するための小文字定義。
	mutex  sync.Mutex
	lv     level.Level
	writer *syslog.Writer
}

func (hndl *SyslogHandler) Output(depth int, lv level.Level, v ...interface{}) {
	hndl.mutex.Lock()
	defer hndl.mutex.Unlock()

	if lv > hndl.lv {
		return
	}

	hndl.mutex.Unlock()
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(depth + 1)
	if !ok {
		file = "???"
		line = 0
	}

	file = trimPrefix(file)

	msg := fmt.Sprintf("%.3v %s:%d %s\n", lv, file, line, fmt.Sprint(v...))

	hndl.mutex.Lock()

	switch lv {
	case level.ERR:
		hndl.writer.Err(msg)
	case level.WARN:
		hndl.writer.Warning(msg)
	case level.INFO:
		hndl.writer.Info(msg)
	case level.DEBUG:
		hndl.writer.Debug(msg)
	}
}

func (hndl *SyslogHandler) SetLevel(lv level.Level) {
	hndl.mutex.Lock()
	defer hndl.mutex.Unlock()

	hndl.lv = lv
}

func (hndl *SyslogHandler) Flush() {
	// hndl.writer.Close()
	return
}

func (hndl *SyslogHandler) Close() error {
	return erro.Wrap(hndl.writer.Close())
}

func NewSyslogHandler(tag string) (Handler, error) {
	writer, err := syslog.New(syslog.LOG_INFO, tag)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &SyslogHandler{writer: writer}, nil
}
