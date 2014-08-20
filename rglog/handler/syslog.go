package handler

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"log/syslog"
)

type syslogCoreHandler struct {
	*syslog.Writer
}

func (hndl *syslogCoreHandler) output(file string, line int, lv level.Level, v ...interface{}) {
	// {レベル} {ファイル名}:{行番号} {メッセージ}
	// 日時は syslog が付ける。
	msg := fmt.Sprintf("%.3v %s:%d %s\n", lv, file, line, fmt.Sprint(v...))

	switch lv {
	case level.ERR:
		hndl.Err(msg)
	case level.WARN:
		hndl.Warning(msg)
	case level.INFO:
		hndl.Info(msg)
	case level.DEBUG:
		hndl.Debug(msg)
	}
}

func (hndl *syslogCoreHandler) flush() {
	return
}

func NewSyslogHandler(tag string) (Handler, error) {
	conn, err := syslog.New(syslog.LOG_INFO, tag)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return wrapCoreHandler(newSynchronizedCoreHandler(&syslogCoreHandler{conn})), nil
}

func NewSyslogHandlerOf(tag, prot, addr string) (Handler, error) {
	conn, err := syslog.Dial(prot, addr, syslog.LOG_INFO, tag)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return wrapCoreHandler(newSynchronizedCoreHandler(&syslogCoreHandler{conn})), nil
}
