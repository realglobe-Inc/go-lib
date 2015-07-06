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
	"sync"

	"github.com/realglobe-Inc/go-lib/rglog/level"
)

// 何もしないハンドラ。
// デバッグ用。
type nopHandler struct {
	lock sync.Mutex
	lv   level.Level
}

func NewNopHandler() Handler {
	return &nopHandler{lv: level.ALL}
}

func (hndl *nopHandler) Level() level.Level {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	return hndl.lv
}

func (hndl *nopHandler) SetLevel(lv level.Level) {
	hndl.lock.Lock()
	defer hndl.lock.Unlock()

	hndl.lv = lv
}

func (hndl *nopHandler) Output(rec Record) {}

func (hndl *nopHandler) Flush() {}

func (hndl *nopHandler) Close() {}
