package erro

import (
	"fmt"
	"runtime"
)

// デバッグ用エラー。
type Tracer struct {
	cause error
	trace string
}

func (err *Tracer) Error() string {
	return err.cause.Error() + "\n" +
		//"--------------------------------------------------\n" +
		err.trace
}

func (err *Tracer) Cause() error {
	return err.cause
}

const traceLen = 8192

// スタックトレースを付加する。
// nil をそのまま返すので、 return Wrap(func() error) みたいな使い方ができる。
// 囲うのは最初だけでも良いし、既に Wrap されている場合はそのまま返すので、
// 最初かどうか分からなければ、毎回囲っても良い。
// スタックトレースの先頭はこの関数になるが、気にしない。
func Wrap(err error) error {
	if err == nil {
		return nil
	} else if _, ok := err.(*Tracer); ok {
		return err
	}

	buff := make([]byte, traceLen)
	len := runtime.Stack(buff, false)

	for len > 0 && buff[len-1] == '\n' {
		len--
	}

	return &Tracer{err, string(buff[:len])}
}

func Unwrap(err error) error {
	for {
		tr, ok := err.(*Tracer)
		if ok {
			err = tr.cause
		} else {
			return err
		}
	}
}

type Error string

func (err Error) Error() string {
	return string(err)
}

// 引数から 1 つの文字列をつくって最初のエラーにして、それを Wrap して返す。
// 文字列の成形は fmt.Print() 形式。
func New(a ...interface{}) error {
	return Wrap(Error(fmt.Sprint(a...)))
}
