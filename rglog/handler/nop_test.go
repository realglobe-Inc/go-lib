package handler

import (
	"testing"
)

func TestNopHandlerLevel(t *testing.T) {
	testHandlerLevel(t, NewNopHandler())
}

func TestNopHandlerOutput(t *testing.T) {
	testHandlerOutput(t, NewNopHandler())
}

func BenchmarkNopHandler(b *testing.B) {
	benchmarkHandler(b, NewNopHandler())
}
