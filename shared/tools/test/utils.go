package test

import (
	"testing"
)

// Ok asserts no error
func Ok(t *testing.T, e error) {
	if e != nil {
		t.Fatal(e)
	}
}

// Assert expression
func Assert(t *testing.T, cond bool, msg string) {
	if !cond {
		t.Fatal(msg)
	}
}

// Assertf condition
func Assertf(t *testing.T, cond bool, f string, v ...interface{}) {
	if !cond {
		t.Fatalf(f, v...)
	}
}

// AssertExpAct condition
func AssertExpAct(t *testing.T, exp interface{}, act interface{}) {
	Assertf(t, exp == act, "Wrong value:\nwant: %v\ngot: %v\n", exp, act)
}
