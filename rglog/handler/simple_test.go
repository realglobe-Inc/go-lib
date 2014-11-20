package handler

import (
	"io/ioutil"
	"testing"
)

func TestBasicHandlerLevel(t *testing.T) {
	testHandlerLevel(t, NewBasicHandler(ioutil.Discard))
}

func TestBasicHandlerOutput(t *testing.T) {
	testHandlerOutput(t, NewBasicHandler(ioutil.Discard))
}

func BenchmarkBasicHandler(b *testing.B) {
	benchmarkHandler(b, NewBasicHandler(ioutil.Discard))
}

func TestFlushHandlerLevel(t *testing.T) {
	testHandlerLevel(t, NewFlushHandler(ioutil.Discard))
}

func TestFlushHandlerOutput(t *testing.T) {
	testHandlerOutput(t, NewFlushHandler(ioutil.Discard))
}

func BenchmarkFlushHandler(b *testing.B) {
	benchmarkHandler(b, NewFlushHandler(ioutil.Discard))
}
