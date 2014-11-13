package handler

import (
	"testing"
)

func TestNopHandler(t *testing.T) {
	testHandler(t, NewNopHandler())
}

func BenchmarkNopHandler(b *testing.B) {
	benchmarkHandler(b, NewNopHandler())
}
