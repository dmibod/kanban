package update_test

import (
	"github.com/stretchr/testify/mock"
	"context"
	"github.com/go-chi/chi"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmibod/kanban/shared/tools/logger/noop"

	"github.com/dmibod/kanban/shared/services"

	_service "github.com/dmibod/kanban/shared/services/mocks"

	_factory "github.com/dmibod/kanban/update/mocks"

	"github.com/dmibod/kanban/shared/kernel"

	"github.com/dmibod/kanban/update"
)

func TestCreateCard(t *testing.T) {

	payload := &update.Card{ID: "5c16dd24c7ee6e5dcf626266", Name: "Sample"}

	model := &services.CardPayload{Name: payload.Name}

	service := &_service.CardService{}
	service.On("CreateCard", model).Return(kernel.Id(payload.ID), nil).Once()

	req := toJsonRequest(t, http.MethodPost, "http://localhost/v1/api/card/", payload)
	res := httptest.NewRecorder()

	getAPI(service).Create(res, req)

	service.AssertExpectations(t)

	expected := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{payload.ID, true}

	exp := strings.TrimSpace(string(toJson(t, &expected)))
	act := strings.TrimSpace(res.Body.String())

	assertf(t, act == exp, "Wrong response\nwant: %v\ngot: %v", exp, act)
}

func getAPI(s services.CardService) *update.API {
	factory := &_factory.ServiceFactory{}
	factory.On("CreateCardService", mock.Anything).Return(s)
	return update.CreateAPI(&noop.Logger{}, factory)
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

func toJsonRequest(t *testing.T, m string, u string, o interface{}, f ...func(*chi.Context)) *http.Request {
	r, err := http.NewRequest(m, u, bytes.NewBuffer(toJson(t, o)))
	ok(t, err)
	return toChiRequest(r, f...)
}

func toRequest(t *testing.T, m string, u string, f ...func(*chi.Context)) *http.Request {
	r, err := http.NewRequest(m, u, bytes.NewBuffer([]byte{}))
	ok(t, err)
	return toChiRequest(r, f...)
}

func toChiRequest(r *http.Request, f ...func(*chi.Context)) *http.Request {
	rctx := chi.NewRouteContext() 
	for _, i := range f {
		i(rctx)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
