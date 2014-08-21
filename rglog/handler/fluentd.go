package handler

import (
	"bufio"
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"net"
	"os"
	"strconv"
	"time"
)

// fluentd の in_forward にログを流す coreHandler。
type fluentdCoreHandler struct {
	// fluentd の tag。
	tag string

	// fluentd サーバーへの接続端。
	conn net.Conn
	*bufio.Writer
}

func (hndl *fluentdCoreHandler) output(file string, line int, lv level.Level, v ...interface{}) {
	// 形式は JSON で書けば、
	//
	// [
	//   "a.b.c",
	//   1308466941,
	//   {
	//     "level": "INFO",
	//     "file": "github.com/realglobe-Inc/go-lib-rg/rglog/handler",
	//     "line": 39,
	//     "message": "Unko!"
	//   }
	// ]
	//
	// これを MessagePack にして送る。
	// http://docs.fluentd.org/ja/articles/in_forward と
	// https://github.com/msgpack/msgpack/blob/master/spec.md を参照。

	date := time.Now().Unix()
	msg := fmt.Sprint(v...)

	// fixarray 3.
	buff := []byte{0x90 | 3}

	buff = append(buff, messagePackString(hndl.tag)...)

	buff = append(buff, messagePackInteger(date)...)

	// fixmap 4.
	buff = append(buff, 0x80|4)

	buff = append(buff, messagePackString("level")...)
	buff = append(buff, messagePackString(lv.String())...)

	buff = append(buff, messagePackString("file")...)
	buff = append(buff, messagePackString(file)...)

	buff = append(buff, messagePackString("line")...)
	buff = append(buff, messagePackInteger(int64(line))...)

	buff = append(buff, messagePackString("message")...)
	buff = append(buff, messagePackString(msg)...)

	if _, err := hndl.Write(buff); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
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

func (hndl *fluentdCoreHandler) flush() {
	if err := hndl.Flush(); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
}

func (hndl *fluentdCoreHandler) close() {
	hndl.flush()

	if err := hndl.conn.Close(); err != nil {
		err = erro.Wrap(err)
		fmt.Fprintln(os.Stderr, err)
	}
}

func NewFluentdHandler(addr, tag string) (Handler, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	return wrapCoreHandler(newSynchronizedCoreHandler(&fluentdCoreHandler{tag, conn, bufio.NewWriter(conn)})), nil
}
