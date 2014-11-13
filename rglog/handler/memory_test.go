package handler

import (
	"testing"
)

func TestMemoryHandler(t *testing.T) {
	testHandler(t, NewMemoryHandler())
}

func BenchmarkMemoryHandler(b *testing.B) {
	benchmarkHandler(b, NewMemoryHandler())
}
