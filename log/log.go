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
	DEBUG
)

func (level Level) String() string {
	switch level {
	case ERR:
		return "ERR"
	case INFO:
		return "INFO"
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

func Debug(v ...interface{}) {
	logging(DEBUG, v...)
}
