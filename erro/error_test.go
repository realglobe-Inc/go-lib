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

package erro

import (
	"errors"
	"reflect"
	"testing"
)

func TestWrapNil(t *testing.T) {
	if Wrap(nil) != nil {
		t.Error("not nil")
	}
}

func TestWrapTracer(t *testing.T) {
	tr := New("test")
	if Wrap(tr) != tr {
		t.Error("not through")
	}
}

func TestWrap(t *testing.T) {
	msg := "test"
	err := New(msg)

	if tr, ok := err.(*Tracer); !ok {
		t.Error(reflect.TypeOf(err))
	} else if tr.Stack() == "" {
		t.Error(tr.Stack())
	} else if cause := tr.Cause(); cause.Error() != msg {
		t.Error(cause.Error(), msg)
	} else if m := tr.Error(); len(m) <= len(cause.Error()) || len(m) <= len(tr.Stack()) {
		t.Error(m)
		t.Error(cause.Error())
		t.Error(tr.Stack())
	}
}

func TestUnwrapNil(t *testing.T) {
	if Unwrap(nil) != nil {
		t.Error("not nil")
	}
}

func TestUnwrapNonTracer(t *testing.T) {
	err := errors.New("test")
	if Unwrap(err) != err {
		t.Error("not through")
	}
}

func TestUnwrap(t *testing.T) {
	msg := "test"
	err := New(msg)

	if cause := Unwrap(err); cause == err {
		t.Error(cause)
	} else if _, ok := cause.(*Tracer); ok {
		t.Error(cause)
	} else if cause.Error() != msg {
		t.Error(cause.Error(), msg)
	}
}
