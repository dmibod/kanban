package update_test

import (
	"bytes"
	"net/http"
	"encoding/json"
	"testing"
)

func ok(t *testing.T, e error) {
	if e != nil {
		t.Fatal(e)
	}
}

func assert(t *testing.T, exp bool, msg string) {
	if !exp {
		t.Fatal(msg)
	}
}

func assertf(t *testing.T, exp bool, f string, v ...interface{}) {
	if !exp {
		t.Fatalf(f, v...)
	}
}

func toJson(t *testing.T, o interface{}) []byte {
	bytes, err := json.Marshal(o)
	ok(t, err)
	return bytes
}

func toJsonRequest(t *testing.T, m string, u string, o interface{}) *http.Request {
	r, err := http.NewRequest(m, u, bytes.NewBuffer(toJson(t, o)))
	ok(t, err)
	return r
}
