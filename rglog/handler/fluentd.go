package handler

import (
	"bytes"
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// fluentd の http 入力にログを流すハンドラ。

const bodyType = "application/x-www-form-urlencoded"

type FluentdHandler struct {
	lock   sync.Mutex
	lv     level.Level
	url    string
	client *http.Client
}

var notPrintable *regexp.Regexp

func init() {
	notPrintable = regexp.MustCompile("[^[:print:]]+")
}

func (hndl *FluentdHandler) Output(depth int, lv level.Level, v ...interface{}) {
	hndl.lock.Lock()
	if lv > hndl.lv {
		return
	}
	hndl.lock.Unlock()

	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(depth + 1)
	if !ok {
		file = "???"
		line = 0
	}

	file = trimPrefix(file)

	// {"level":"{レベル}","file":"{ファイル名}","line":{行番号},"message":{メッセージ}}
	// 日時は fluentd が付ける。
	msg := fmt.Sprint(v...)
	msg = strings.Replace(msg, "\"", "\\\"", -1)
	msg = notPrintable.ReplaceAllString(msg, " ")
	buff := fmt.Sprintf("json={\"level\":\"%v\",\"file\":\"%s\",\"line\":%d,\"message\":\"%s\"}", lv, file, line, msg)

	// http.Client がスレッドセーフらしいのでロックしなくていい。
	// TODO fluentd が処理して返答するまで待ってしまうので非効率。
	// fluentd からの返答を捨て続けるゴルーチンを立てて、送信したら後は知らんという形に書き換えた方が良さそう。
	resp, err := hndl.client.Post(hndl.url, bodyType, bytes.NewBufferString(buff))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if err := resp.Body.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (hndl *FluentdHandler) SetLevel(lv level.Level) {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	hndl.lv = lv
}

func (hndl *FluentdHandler) Flush() {
	return
}

func NewFluentdHandler(host, tag string) (Handler, error) {
	// 接続テスト。
	client := &http.Client{}
	resp, err := client.Post("http://"+host+"/debug.test", bodyType, bytes.NewBufferString("json={\"debug\":\"test\"}"))
	if err != nil {
		return nil, erro.Wrap(err)
	}
	if err := resp.Body.Close(); err != nil {
		return nil, erro.Wrap(err)
	}

	return &FluentdHandler{url: "http://" + host + "/" + tag, client: client}, nil
}
