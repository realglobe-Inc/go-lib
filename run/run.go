package run

import (
	"bytes"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"os"
	"os/exec"
	"strings"
)

// 会話型。
func Run(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debug(cmd.Args)

	return erro.Wrap(cmd.Run())
}

// 非会話型。
func NonInteractive(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debug(cmd.Args)

	return erro.Wrap(cmd.Run())
}

// 非会話型でエラーを無視する。
func Neglect(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debug(cmd.Args)

	if err := cmd.Run(); err != nil {
		log.Err(erro.Unwrap(err))
		log.Debug(err)
	}
}

type Error struct {
	cause  error
	stdout string
	stderr string
}

func (err *Error) Error() string {
	return err.cause.Error() + " Stdout[" + err.stdout + "] Stderr[" + err.stderr + "]"
}

func (err *Error) Cause() error {
	return err.cause
}

func (err *Error) Stdout() string {
	return err.stdout
}

func (err *Error) Stderr() string {
	return err.stderr
}

func newError(cause error, stdout, stderr string) *Error {
	return &Error{cause, strings.TrimSpace(stdout), strings.TrimSpace(stderr)}
}

// 画面表示せず、非会話型で画面出力を返す。
func Output(args ...string) (string, string, error) {
	cmd := exec.Command(args[0], args[1:]...)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	log.Debug(cmd.Args)

	if err := cmd.Run(); err != nil {
		return "", "", erro.Wrap(newError(err, stdout.String(), stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), nil
}

// 画面表示せず、非会話型で標準出力を返す。
func Stdout(args ...string) (string, error) {
	stdout, _, err := Output(args...)
	if err != nil {
		return "", err
	}

	return stdout, nil
}

// 画面表示せず、非会話型。
func Quiet(args ...string) error {
	if _, _, err := Output(args...); err != nil {
		return err
	}

	return nil
}
