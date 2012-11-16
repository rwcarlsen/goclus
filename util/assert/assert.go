package assert

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

type T struct{ *testing.T }

func (t *T) Fatal() {
	if t.Failed() {
		t.FailNow()
	}
}

func Eq(t *testing.T, i, j interface{}) *T {
	if i != j {
		t.Errorf("%v %v != %v", caller(), i, j)
	}
	return &T{t}
}

func Ne(t *testing.T, i, j interface{}) *T {
	if i == j {
		t.Errorf("%v %v == %v", caller(), i, j)
	}
	return &T{t}
}

func NoErr(t *testing.T, err error) *T {
	if err != nil {
		t.Error(caller(), " error: ", err)
	}
	return &T{t}
}

func Err(t *testing.T, err error) *T {
	if err == nil {
		t.Error(caller(), " expected error, got nil")
	}
	return &T{t}
}

func caller() string {
	_, file, line, ok := runtime.Caller(2)
	msg := ""
	if ok {
		msg = fmt.Sprint(filepath.Base(file), ":", line)
	}
	return msg
}
