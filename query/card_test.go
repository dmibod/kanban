package query_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/stretchr/testify/mock"

	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/services/mocks"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

func TestGetCard(t *testing.T) {

	id := "5c16dd24c7ee6e5dcf626266"

	model := &services.CardModel{ID: kernel.Id(id), Name: "Sample"}

	service := &mocks.CardService{}
	service.On("GetCardByID", mock.Anything, kernel.Id(id)).Return(model, nil).Once()

	req := toRequest(t, http.MethodGet, "http://localhost/v1/api/card/"+id, func(rctx *chi.Context) {
		rctx.URLParams.Add("ID", id)
	})

	rec := httptest.NewRecorder()

	getAPI(service).Get(rec, req)

	res := rec.Result()

	service.AssertExpectations(t)

	expected := &query.Card{
		ID:   string(model.ID),
		Name: model.Name,
	}

	exp := strings.TrimSpace(string(toJson(t, expected)))
	act := strings.TrimSpace(string(body(t, res)))
	test.AssertExpAct(t, exp, act)
}

func getAPI(s services.CardService) *query.CardAPI {
	return query.CreateCardAPI(&noop.Logger{}, s)
}

func body(t *testing.T, res *http.Response) []byte {
	body, err := ioutil.ReadAll(res.Body)
	test.Ok(t, err)
	return body
}

func toJson(t *testing.T, o interface{}) []byte {
	bytes, err := json.Marshal(o)
	test.Ok(t, err)
	return bytes
}

func toJsonRequest(t *testing.T, m string, u string, o interface{}, f ...func(*chi.Context)) *http.Request {
	r, err := http.NewRequest(m, u, bytes.NewBuffer(toJson(t, o)))
	test.Ok(t, err)
	return toChiRequest(r, f...)
}

func toRequest(t *testing.T, m string, u string, f ...func(*chi.Context)) *http.Request {
	r, err := http.NewRequest(m, u, bytes.NewBuffer([]byte{}))
	test.Ok(t, err)
	return toChiRequest(r, f...)
}

func toChiRequest(r *http.Request, f ...func(*chi.Context)) *http.Request {
	rctx := chi.NewRouteContext()
	for _, i := range f {
		i(rctx)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
