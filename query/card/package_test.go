package card_test

import (
	"github.com/dmibod/kanban/shared/services/card/mocks"
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmibod/kanban/query/card"

	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/stretchr/testify/mock"

	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/kernel"
	service "github.com/dmibod/kanban/shared/services/card"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

func TestGetCard(t *testing.T) {

	id := "5c16dd24c7ee6e5dcf626266"

	model := &service.Model{ID: kernel.ID(id), Name: "Sample"}

	service := &mocks.Service{}
	service.On("GetByID", mock.Anything, kernel.ID(id)).Return(model, nil).Once()

	req := toRequest(t, http.MethodGet, "http://localhost/v1/api/card/"+id, func(rctx *chi.Context) {
		rctx.URLParams.Add("CARDID", id)
	})

	rec := httptest.NewRecorder()

	getAPI(service).Get(rec, req)

	res := rec.Result()

	service.AssertExpectations(t)

	expected := &card.Card{
		ID:   string(model.ID),
		Name: model.Name,
	}

	exp := strings.TrimSpace(string(toJson(t, expected)))
	act := strings.TrimSpace(string(body(t, res)))
	test.AssertExpAct(t, exp, act)
}

func getAPI(s service.Service) *card.API {
	return card.CreateAPI(s, &noop.Logger{})
}

func body(t *testing.T, res *http.Response) []byte {
	body, err := ioutil.ReadAll(res.Body)
	test.Ok(t, err)
	return body
}

func toJson(t *testing.T, payload interface{}) []byte {
	bytes, err := json.Marshal(payload)
	test.Ok(t, err)
	return bytes
}

func toJsonRequest(t *testing.T, method string, url string, payload interface{}, f ...func(*chi.Context)) *http.Request {
	r, err := http.NewRequest(method, url, bytes.NewBuffer(toJson(t, payload)))
	test.Ok(t, err)
	return toChiRequest(r, f...)
}

func toRequest(t *testing.T, method string, url string, f ...func(*chi.Context)) *http.Request {
	r, err := http.NewRequest(method, url, bytes.NewBuffer([]byte{}))
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
