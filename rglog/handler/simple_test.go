// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
