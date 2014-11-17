package handler

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"log/syslog"
	"os"
	"strconv"
)

// syslog にログを流す coreHandler。
// ログデーモンが一時的に落ちていても、動き出せば元通りに動く。
type syslogCoreHandler struct {
	tag  string
	addr string

	base *syslog.Writer
}

func (core *syslogCoreHandler) output(rec Record) {
	// {レベル} {ファイル名}:{行番号} {メッセージ}
	// 日時は syslog が付ける。
	msg := fmt.Sprintf("%."+strconv.Itoa(lvWidth)+"v %s:%d %s\n",
		rec.Level(), rec.File(), rec.Line(), rec.Message())

	for retry := false; ; retry = true {
		if core.base == nil {
			var err error
			core.base, err = syslog.Dial("", core.addr, syslog.LOG_INFO, core.tag)
			if err != nil {
				// 初期化出来なければ諦める。
				fmt.Fprintln(os.Stderr, erro.Wrap(err))
				fmt.Fprintln(os.Stderr, "Drop log: "+string(SimpleFormatter.Format(rec)))
				return
			}
		}

		// 初期化してある。

		var err error
		switch rec.Level() {
		case level.ERR:
			err = core.base.Err(msg)
		case level.WARN:
			err = core.base.Warning(msg)
		case level.INFO:
			err = core.base.Info(msg)
		case level.DEBUG:
			err = core.base.Debug(msg)
		}
		if err == nil {
			// 書き込み成功。
			return
		}

		// 書き込み失敗。
		// 初期化が古くてサーバー側で何か変わったとか。

		fmt.Fprintln(os.Stderr, erro.Wrap(err))
		core.base.Close()
		core.base = nil

		if retry {
			// 初期化しなおしても書き込めないなら諦める。
			fmt.Fprintln(os.Stderr, "Drop log: "+string(SimpleFormatter.Format(rec)))
			return
		}
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
	return NewSyslogHandlerTo("", tag)
}

func NewSyslogHandlerTo(addr, tag string) Handler {
	return wrapCoreHandler(newSynchronizedCoreHandler(&syslogCoreHandler{tag: tag, addr: addr}))
}
