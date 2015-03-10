// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"container/list"
	"fmt"
	"github.com/realglobe-Inc/go-lib/erro"
	"io"
	"net"
	"os"
	"strconv"
)

// fluentd の in_forward にログを流す coreHandler。
// fluentd が一時的に落ちていても、動き出せば元通りに動く。
// ただし、落ちていた間のログは欠落する。
type fluentdCoreHandler struct {
	// fluentd の tag。
	tag  string
	addr string

	// fluentd への接続。
	conn net.Conn
	// バッファ兼書き込み口。
	buff *tcpLogBuffer
}

func (core *fluentdCoreHandler) output(rec Record) {
	// 形式は JSON で書けば、
	//
	// [
	//   "a.b.c",
	//   1308466941,
	//   {
	//     "level": "INFO",
	//     "file": "github.com/realglobe-Inc/go-lib/rglog/handler",
	//     "line": 39,
	//     "message": "Unko!"
	//   }
	// ]
	//
	// これを MessagePack にして送る。
	// http://docs.fluentd.org/ja/articles/in_forward と
	// https://github.com/msgpack/msgpack/blob/master/spec.md を参照。

	// fixarray 3.
	buff := []byte{0x90 | 3}

	buff = append(buff, messagePackString(core.tag)...)

	buff = append(buff, messagePackInteger(rec.Date().Unix())...)

	// fixmap 4.
	buff = append(buff, 0x80|4)

	buff = append(buff, messagePackString("level")...)
	buff = append(buff, messagePackString(rec.Level().String())...)

	buff = append(buff, messagePackString("file")...)
	buff = append(buff, messagePackString(rec.File())...)

	buff = append(buff, messagePackString("line")...)
	buff = append(buff, messagePackInteger(int64(rec.Line()))...)

	buff = append(buff, messagePackString("message")...)
	buff = append(buff, messagePackString(rec.Message())...)

	const ( // てきとう。
		writeSize  = 4096
		bufferSize = 2*writeSize + 1024
	)

	core.buff.add(buff)
	defer core.buff.trim(bufferSize)
	core.flushCore(writeSize)
}

func messagePackString(val string) []byte {
	buff := []byte{}
	if length := len(val); length < (1 << 4) {
		// fixstr.
		buff = append(buff, byte(0xa0|length))
	} else if length < (1 << 8) {
		// str8.
		buff = append(buff, 0xd9, byte(length))
	} else if length < (1 << 16) {
		// str16.
		buff = append(buff, 0xda, byte((length&(0xff<<8))>>8), byte(length&0xff))
	} else if length < (1 << 32) {
		// str32.
		buff = append(buff, 0xdb, byte((length&(0xff<<24))>>24), byte((length&(0xff<<16))>>16), byte((length&(0xff<<8))>>8), byte(length&0xff))
	} else {
		panic("too long string " + strconv.Itoa(length) + ".")
	}
	buff = append(buff, []byte(val)...)
	return buff
}

func messagePackInteger(val int64) []byte {
	buff := []byte{}
	if -(1<<7) <= val && val < (1<<7) {
		// int8.
		buff = append(buff, 0xd0, byte(val))
	} else if -(1<<15) <= val && val < (1<<15) {
		// int16.
		buff = append(buff, 0xd1, byte((val&(0xff<<8))>>8), byte(val&0xff))
	} else if -(1<<31) <= val && val < (1<<31) {
		// int32.
		buff = append(buff, 0xd2, byte((val&(0xff<<24))>>24), byte((val&(0xff<<16))>>16), byte((val&(0xff<<8))>>8), byte(val&0xff))
	} else {
		// int64.
		buff = append(buff, 0xd3, byte((uint64(val)&uint64(0xff<<56))>>56), byte((val&(0xff<<48))>>48), byte((val&(0xff<<40))>>40), byte((val&(0xff<<32))>>32), byte((val&(0xff<<24))>>24), byte((val&(0xff<<16))>>16), byte((val&(0xff<<8))>>8), byte(val&0xff))
	}
	return buff
}

func (core *fluentdCoreHandler) flush() {
	core.flushCore(0)
}

func (core *fluentdCoreHandler) flushCore(writeSize int) {
	for retry := false; ; retry = true {
		if core.conn == nil {
			var err error
			core.conn, err = net.Dial("tcp", core.addr)
			if err != nil {
				// 接続出来なければ諦める。
				fmt.Fprintln(os.Stderr, erro.Wrap(err))
				return
			}
			core.buff.setSink(core.conn)
		}

		// 接続はある。

		err := core.buff.flushIfNeeded(writeSize)
		if err == nil {
			// 書き込み成功。
			return
		}

		// 書き込み失敗。
		// 接続が古くてサーバー側に切断されていたとか。

		fmt.Fprintln(os.Stderr, erro.Wrap(err))
		core.conn.Close()
		core.conn = nil

		if retry {
			// 接続しなおしても書き込めないなら諦める。
			return
		}
	}
}

func (core *fluentdCoreHandler) close() {
	if core.conn == nil {
		return
	}
	core.flush()
	core.buff.close()
	core.conn.Close()
	core.conn = nil
}

func NewFluentdHandler(addr, tag string) Handler {
	return wrapCoreHandler(newSynchronizedCoreHandler(&fluentdCoreHandler{tag: tag, addr: addr, buff: newTcpLogBuffer()}))
}

// TCP 接続用のバッファ付き書き込みラッパー。
// TCP 接続はサーバー側が閉じた後でも、1 回は書き込みエラーが出ないので、
// その場合にもログを消失させずに、再接続後に送り直せるようにする。
type tcpLogBuffer struct {
	// 書き込み口。ユーザー空間でバッファリングしていない TCP 接続を想定。
	sink io.Writer

	// 前回の sink.Write 時に書き出した分。先頭が一番古いキュー。
	written     *list.List
	writtenSize int
	// まだ書き込んでない分。先頭が一番古いキュー。
	unwritten     *list.List
	unwrittenSize int
}

func newTcpLogBuffer() *tcpLogBuffer {
	return &tcpLogBuffer{
		written:   list.New(),
		unwritten: list.New(),
	}
}

func (buff *tcpLogBuffer) setSink(sink io.Writer) {
	buff.sink = sink

	// 前の接続に書き出せていない可能性のあるログを移す。
	buff.unwritten.PushFrontList(buff.written)
	buff.unwrittenSize += buff.writtenSize
	buff.written = list.New()
	buff.writtenSize = 0
}

func (buff *tcpLogBuffer) add(log []byte) {
	buff.unwritten.PushBack(log)
	buff.unwrittenSize += len(log)
}

// サイズが baseline 未満にならないところまで古いログを消す。
func (buff *tcpLogBuffer) trim(baseline int) {
	for buff.writtenSize+buff.unwrittenSize > baseline && buff.written.Len() > 0 {
		elem := buff.written.Front()
		log := elem.Value.([]byte)
		if buff.writtenSize-len(log)+buff.unwrittenSize < baseline {
			return
		}
		buff.written.Remove(elem)
		buff.writtenSize -= len(log)
		fmt.Fprintln(os.Stderr, "Maybe drop log: "+string(log))
	}
	for buff.unwrittenSize > baseline && buff.unwritten.Len() > 0 {
		elem := buff.unwritten.Front()
		log := elem.Value.([]byte)
		if buff.unwrittenSize-len(log) < baseline {
			return
		}
		buff.unwritten.Remove(elem)
		buff.unwrittenSize -= len(log)
		fmt.Fprintln(os.Stderr, "Drop log: "+string(log))
	}
}

func (buff *tcpLogBuffer) flushIfNeeded(unwrittenLimit int) error {
	if buff.unwrittenSize <= unwrittenLimit {
		return nil
	}
	return buff.flush()
}

func (buff *tcpLogBuffer) flush() error {
	writeBuff := []byte{}
	for elem := buff.unwritten.Front(); elem != nil; elem = elem.Next() {
		log := elem.Value.([]byte)
		writeBuff = append(writeBuff, log...)
	}
	if _, err := buff.sink.Write(writeBuff); err != nil {
		return erro.Wrap(err)
	}

	// 前回の分は書き出せてた。

	buff.written = buff.unwritten
	buff.writtenSize = buff.unwrittenSize
	buff.unwritten = list.New()
	buff.unwrittenSize = 0
	return nil
}

func (buff *tcpLogBuffer) close() {
	buff.flush()
	*buff = tcpLogBuffer{}
}
