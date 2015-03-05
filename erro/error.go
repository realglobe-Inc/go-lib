package erro

import (
	"errors"
	"fmt"
	"runtime"
)

// スタックトレース付きエラー。
type Tracer struct {
	cause error
	trace string
}

// error を実装。
func (tr *Tracer) Error() string {
	return tr.Cause().Error() + "\n" + tr.Stack()
}

// 素のエラーを返す。
func (tr *Tracer) Cause() error {
	return tr.cause
}

// 表示用スタックトレースを返す。
func (tr *Tracer) Stack() string {
	return tr.trace
}

// 表示用スタックトレースの最大バイト長。
const traceLen = 8192

// スタックトレースを付加する。スタックトレースの先頭はこの関数になってしまうが、気にするな。
// nil はそのまま返すので、 return Wrap(func() error) みたいな使い方もできる。
// 既に Wrap されている場合はそのまま返すので、毎回 Wrap しても良い。
func Wrap(err error) error {
	if err == nil {
		return nil
	} else if tr, ok := err.(*Tracer); ok {
		return tr
	}

	buff := make([]byte, traceLen)
	runtime.Stack(buff, false)

	// 普通、Error() の返り値の末尾に改行は付かないので、末尾の改行を削除する。
	for ; len(buff) > 0 && buff[len(buff)-1] == '\n'; buff = buff[:len(buff)-1] {
	}

	return &Tracer{err, string(buff)}
}

// スタックトレース付きエラーだったら、素のエラーを取り出す。
// そうでなければ、そのまま返す。
func Unwrap(err error) error {
	if err == nil {
		return nil
	} else if tr, ok := err.(*Tracer); ok {
		return tr.Cause()
	} else {
		return err
	}
}

// スタックトレース付きのエラーをつくる。
// 素のエラーは erros.New(fmt.Sprint(a...)) でつくる。
func New(a ...interface{}) error {
	return Wrap(errors.New(fmt.Sprint(a...)))
}
