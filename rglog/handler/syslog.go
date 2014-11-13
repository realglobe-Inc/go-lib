package handler

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"log/syslog"
	"os"
)

type syslogCoreHandler struct {
	tag string

	base *syslog.Writer
}

func (core *syslogCoreHandler) output(file string, line int, lv level.Level, v ...interface{}) {
	// {レベル} {ファイル名}:{行番号} {メッセージ}
	// 日時は syslog が付ける。
	msg := fmt.Sprintf("%.3v %s:%d %s\n", lv, file, line, fmt.Sprint(v...))

	// ログデーモンが一時的に落ちていても、動き出せば元通りに動くように。
	if core.base == nil {
		var err error
		core.base, err = syslog.New(syslog.LOG_INFO, core.tag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}

	var err error
	switch lv {
	case level.ERR:
		err = core.base.Err(msg)
	case level.WARN:
		err = core.base.Warning(msg)
	case level.INFO:
		err = core.base.Info(msg)
	case level.DEBUG:
		err = core.base.Debug(msg)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, erro.Wrap(err))
		if err := core.base.Close(); err != nil {
			fmt.Fprintln(os.Stderr, erro.Wrap(err))
		}
		core.base = nil
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

func NewSyslogHandler(tag string) Handler {
	return wrapCoreHandler(newSynchronizedCoreHandler(&syslogCoreHandler{tag: tag}))
}
