package update_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"

	"github.com/dmibod/kanban/shared/tools/logger/noop"

	"github.com/dmibod/kanban/shared/services"

	_service "github.com/dmibod/kanban/shared/services/mocks"

	"github.com/dmibod/kanban/shared/kernel"

	"github.com/dmibod/kanban/update"
)

func TestCardAPI(t *testing.T) {
	id := "5c16dd24c7ee6e5dcf626266"
	testCreateCard(t, id)
	testUpdateCard(t, id)
	testRemoveCard(t, id)
}

func testCreateCard(t *testing.T, id string) {

	payload := &update.Card{ID: id, Name: "Sample"}

	model := &services.CardPayload{Name: payload.Name}

	service := &_service.CardService{}
	service.On("CreateCard", mock.Anything, model).Return(kernel.Id(payload.ID), nil).Once()

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

func testUpdateCard(t *testing.T, id string) {

	model := &services.CardModel{ID: kernel.Id(id), Name: "Sample!"}

	service := &_service.CardService{}
	service.On("UpdateCard", mock.Anything, model).Return(model, nil).Once()

	req := toJsonRequest(t, http.MethodPut, "http://localhost/v1/api/card/"+id, model, func(rctx *chi.Context) {
		rctx.URLParams.Add("ID", id)
	})
	res := httptest.NewRecorder()

	getAPI(service).Update(res, req)

	service.AssertExpectations(t)

	expected := &update.Card{
		ID:   string(model.ID),
		Name: model.Name,
	}

	exp := strings.TrimSpace(string(toJson(t, &expected)))
	act := strings.TrimSpace(res.Body.String())

	assertf(t, act == exp, "Wrong response\nwant: %v\ngot: %v", exp, act)
}

func testRemoveCard(t *testing.T, id string) {

	service := &_service.CardService{}
	service.On("RemoveCard", mock.Anything, kernel.Id(id)).Return(nil).Once()

	req := toRequest(t, http.MethodDelete, "http://localhost/v1/api/card/"+id, func(rctx *chi.Context) {
		rctx.URLParams.Add("ID", id)
	})
	res := httptest.NewRecorder()

	getAPI(service).Remove(res, req)

	service.AssertExpectations(t)

	expected := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{id, true}

	exp := strings.TrimSpace(string(toJson(t, expected)))
	act := strings.TrimSpace(res.Body.String())

	assertf(t, act == exp, "Wrong response\nwant: %v\ngot: %v", exp, act)
}

func getAPI(s services.CardService) *update.CardAPI {
	return update.CreateCardAPI(&noop.Logger{}, s)
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
