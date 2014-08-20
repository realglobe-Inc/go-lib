package handler

import (
	"io/ioutil"
	"testing"
)

func BenchmarkBasicHandler(b *testing.B) {
	benchmarkHandler(b, NewBasicHandler(ioutil.Discard))
}

func BenchmarkFlushHandler(b *testing.B) {
	benchmarkHandler(b, NewFlushHandler(ioutil.Discard))
}
