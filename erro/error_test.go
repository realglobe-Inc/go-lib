package erro

import (
	"testing"
)

func TestWrap(t *testing.T) {
	if Wrap(nil) != nil {
		t.Error(nil)
	}

	msg := "aho"
	err := New(msg)

	if Wrap(err) != err {
		t.Error(err)
	}

	tr, ok := err.(*Tracer)
	if !ok {
		t.Fatal(err)
	}

	var err2 Error
	err2, ok = tr.Cause().(Error)
	if !ok {
		t.Fatal(tr.Cause())
	}

	if string(err2) != msg {
		t.Error(tr.Cause())
	}
}

func TestUnwrap(t *testing.T) {
	if Unwrap(nil) != nil {
		t.Error(nil)
	}

	msg := "aho"
	if Unwrap(Error(msg)) != Error(msg) {
		t.Error(msg)
	}

	err := New(msg)
	if Unwrap(err) != Error(msg) {
		t.Error(err)
	}
}
