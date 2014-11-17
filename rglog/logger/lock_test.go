package logger

import (
	"testing"
)

func TestLockLoggerHandler(t *testing.T) {
	testLoggerHandler(t, NewLockLoggerManager())
}

func TestLockLoggerLevel(t *testing.T) {
	testLoggerLevel(t, NewLockLoggerManager())
}

func TestLockLoggerUseParent(t *testing.T) {
	testLoggerUseParent(t, NewLockLoggerManager())
}

func TestLockLoggerIsLoggable(t *testing.T) {
	testLoggerIsLoggable(t, NewLockLoggerManager())
}

func TestLockLoggerFileName(t *testing.T) {
	testLoggerFileName(t, NewLockLoggerManager())
}

func TestLockLoggerConcurrent(t *testing.T) {
	testLoggerConcurrent(t, NewLockLoggerManager())
}

func BenchmarkLockLoggerConcurrent(b *testing.B) {
	benchmarkLoggerConcurrent(b, NewLockLoggerManager())
}
