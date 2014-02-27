package log

import (
	"bufio"
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"log"
	"os"
	"sync"
)

// 簡易ログ。
// 画面出力とファイル出力が同時にできる。
// レベル指定ができる。
//
// 毎回ロックするので、速くはない。
// 範囲指定できないので、大規模開発には使えない。

type Level int

const (
	ERR Level = iota + 1
	INFO
	WARN
	DEBUG
)

func (level Level) String() string {
	switch level {
	case ERR:
		return "ERR"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case DEBUG:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

var (
	lock          sync.Mutex
	consoleLevel  Level
	consoleLogger *log.Logger

	file       *os.File
	writer     *bufio.Writer
	fileLevel  Level
	fileLogger *log.Logger
)

func init() {
	SetConsole(INFO)
}

func SetConsole(level Level) {
	setConsole(level, "", 0)
}

func setConsole(level Level, prefix string, flag int) {
	lock.Lock()
	defer lock.Unlock()

	if consoleLogger == nil {
		consoleLogger = log.New(os.Stderr, prefix, flag)
	}
	consoleLevel = level
}

func SetFile(level Level, path string) error {
	return setFile(level, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, path)
}

func setFile(level Level, prefix string, flag int, path string) error {
	lock.Lock()
	defer lock.Unlock()

	oldFile := file
	oldWriter := writer

	newFile, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return erro.Wrap(err)
	}
	file = newFile
	writer = bufio.NewWriter(file)
	fileLevel = level
	fileLogger = log.New(writer, prefix, flag)

	if oldFile != nil {
		if e := oldWriter.Flush(); e != nil {
			return erro.Wrap(e)
		}
		if e := oldFile.Close(); e != nil {
			return erro.Wrap(e)
		}
	}

	return nil
}

func CloseFile() error {
	lock.Lock()
	defer lock.Unlock()

	fileLogger = nil

	if file != nil {
		if e := writer.Flush(); e != nil {
			return erro.Wrap(e)
		}
		if e := file.Close(); e != nil {
			return erro.Wrap(e)
		}
	}

	file = nil
	writer = nil

	return nil
}

func logging(level Level, v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()

	var logger *log.Logger

	logger = consoleLogger
	if logger != nil && level <= consoleLevel {
		logger.Output(3, "["+level.String()+"] "+fmt.Sprint(v...)+"\n")
	}

	logger = fileLogger
	if logger != nil && level <= fileLevel {
		logger.SetPrefix("[" + level.String() + "] ")
		logger.Output(3, fmt.Sprint(v...)+"\n")
	}
}

func Err(v ...interface{}) {
	logging(ERR, v...)
}

func Info(v ...interface{}) {
	logging(INFO, v...)
}

func Warn(v ...interface{}) {
	logging(WARN, v...)
}

func Debug(v ...interface{}) {
	logging(DEBUG, v...)
}

// Loggerのインターフェース
type Logger interface {
	Err(v ...interface{})
	Warn(v ...interface{})
	Info(v ...interface{})
	Debug(v ...interface{})
}

type SimpleLogger struct {
}

func (logger SimpleLogger) Err(v ...interface{}) {
	Err(v)
}
func (logger SimpleLogger) Warn(v ...interface{}) {
	Warn(v)
}
func (logger SimpleLogger) Info(v ...interface{}) {
	Info(v)
}
func (logger SimpleLogger) Debug(v ...interface{}) {
	Debug(v)
}

func GetLogger(name string) Logger {
	lf := GetLoggerRegistroy()
	return lf.GetLogger(name)
}

// Loggerを管理するRegistory
type LoggerRegistroy interface {
	// 指定した名前のLoggerを取得する。
	GetLogger(name string) Logger
	// Loggerを追加する。
	AddLogger(name string, factory func()Logger)
}

// LoggerRegistoryの実体
type loggerRegistroy map[string]Logger

// loggerRegistoryのシングルトンインスタンス
var _loggerRegistroy = createLoggerRegistory()

// LoggerRegistroyの初期化
func createLoggerRegistory() loggerRegistroy {
	lr := loggerRegistroy{}
	// TODO ひとまず
	lr.AddLogger("default", func() Logger {
		return &SimpleLogger{}
	})
	return lr
}

// loggerRegistoryの実装
func (lg loggerRegistroy) GetLogger(name string) Logger {
	return lg[name]
}

// loggerRegistoryの実装
func (lg loggerRegistroy) AddLogger(name string, factory func()Logger) {
	lg[name] = factory()
}

// LoggerRegistoryを取得する
func GetLoggerRegistroy() LoggerRegistroy {
	return _loggerRegistroy
}
