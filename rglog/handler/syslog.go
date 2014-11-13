package handler

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"log/syslog"
	"os"
)

type syslogCoreHandler struct {
	base *syslog.Writer
}

func (core *syslogCoreHandler) output(file string, line int, lv level.Level, v ...interface{}) {
	// {レベル} {ファイル名}:{行番号} {メッセージ}
	// 日時は syslog が付ける。
	msg := fmt.Sprintf("%.3v %s:%d %s\n", lv, file, line, fmt.Sprint(v...))

	switch lv {
	case level.ERR:
		core.base.Err(msg)
	case level.WARN:
		core.base.Warning(msg)
	case level.INFO:
		core.base.Info(msg)
	case level.DEBUG:
		core.base.Debug(msg)
	}
}

func (core *syslogCoreHandler) flush() {
	return
}

func (core *syslogCoreHandler) close() {
	core.flush()

	if err := core.base.Close(); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
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
