package handler

import (
	"io/ioutil"
	"testing"
)

func TestBasicHandler(t *testing.T) {
	testHandler(t, NewBasicHandler(ioutil.Discard))
}
func BenchmarkBasicHandler(b *testing.B) {
	benchmarkHandler(b, NewBasicHandler(ioutil.Discard))
}

func TestFlushHandler(t *testing.T) {
	testHandler(t, NewFlushHandler(ioutil.Discard))
}
func BenchmarkFlushHandler(b *testing.B) {
	benchmarkHandler(b, NewFlushHandler(ioutil.Discard))
}
