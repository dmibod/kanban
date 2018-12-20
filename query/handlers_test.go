package query_test

import (
	"context"
	"github.com/go-chi/chi"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmibod/kanban/query"
	_service "github.com/dmibod/kanban/query/mocks"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

func TestGetCard(t *testing.T) {

	id := "5c16dd24c7ee6e5dcf626266"

	model := &services.CardModel{ID: kernel.Id(id), Name: "Sample"}

	service := &_service.CardService{}
	service.On("GetCardByID", kernel.Id(id)).Return(model, nil).Once()

	api := query.CreateAPI(&noop.Logger{}, service)

	req := toRequest(t, http.MethodGet, "http://localhost/v1/api/card/"+id, func(rctx *chi.Context){
		rctx.URLParams.Add("ID", id)
	})
	res := httptest.NewRecorder()

	api.Get(res, req)

	service.AssertExpectations(t)

	expected := &query.Card{
		ID:   string(model.ID),
		Name: model.Name,
	}

	exp := strings.TrimSpace(string(toJson(t, expected)))
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
