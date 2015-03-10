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

func TestLockLoggerLog(t *testing.T) {
	testLoggerLog(t, NewLockLoggerManager())
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
