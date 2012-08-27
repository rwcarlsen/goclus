
package assert

import "testing"

type T struct{*testing.T}

func (t *T) Fatal() {
  if t.Failed() {
    t.FailNow()
  }
}

func Eq(t *testing.T, i, j interface{}) *T {
  if i != j {
    t.Errorf("%v != %v", i, j)
  }
  return &T{t}
}

func Ne(t *testing.T, i, j interface{}) *T {
  if i == j {
    t.Errorf("%v == %v", i, j)
  }
  return &T{t}
}

func NoErr(t *testing.T, err error) *T {
  if err != nil {
    t.Error("error: ", err)
  }
  return &T{t}
}

func Err(t *testing.T, err error) *T {
  if err == nil {
    t.Error("Expected error, got nil")
  }
  return &T{t}
}
