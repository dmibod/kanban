package update_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmibod/kanban/shared/services"

	_service "github.com/dmibod/kanban/update/mocks"

	"github.com/dmibod/kanban/shared/kernel"

	_log "github.com/dmibod/kanban/shared/tools/log/mocks"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/dmibod/kanban/update"
)
func TestCreateCard(t *testing.T) {

	payload := &update.Card{ID: "5c16dd24c7ee6e5dcf626266", Name: "Sample"}

	model := &services.CardPayload{Name: payload.Name}

	service := &_service.CardService{}
	service.On("CreateCard", model).Return(kernel.Id(payload.ID), nil).Once()

	handler := update.CreateCreateCardHandler(&_log.Logger{}, service)

	req := toJsonRequest(t, http.MethodPost, "http://localhost/post", payload)
	res := httptest.NewRecorder()

	mux.Handle(handler).ServeHTTP(res, req)

	service.AssertExpectations(t)

	expected := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{payload.ID, true}

	exp := strings.TrimSpace(string(toJson(t, &expected)))
	act := strings.TrimSpace(res.Body.String())

	assertf(t, act == exp, "Wrong response\nwant: %v\ngot: %v", exp, act)
}

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
