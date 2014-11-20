package handler

import (
	"testing"
)

func TestMemoryHandlerLevel(t *testing.T) {
	testHandlerLevel(t, NewMemoryHandler())
}

func TestMemoryHandlerOutput(t *testing.T) {
	testHandlerOutput(t, NewMemoryHandler())
}

func BenchmarkMemoryHandler(b *testing.B) {
	benchmarkHandler(b, NewMemoryHandler())
}
