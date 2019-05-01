package board_test

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"context"
	"github.com/dmibod/kanban/shared/tools/test"
	"bytes"
	"net/http"
	"github.com/go-chi/chi"
	api "github.com/dmibod/kanban/query/board"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services/board"
	"github.com/stretchr/testify/mock"
	"github.com/dmibod/kanban/shared/services/board/mocks"
	"testing"
	"net/http/httptest"
)

func TestList(t *testing.T) {

	id := "5c16dd24c7ee6e5dcf626266"

	model := &board.ListModel{ID: kernel.ID(id), Name: "Sample"}

	service := &mocks.Service{}
	service.On("GetByOwner", mock.Anything, mock.Anything).Return([]*board.ListModel{model}, nil).Once()

	req := toRequest(t, http.MethodGet, "", func(rctx *chi.Context) {
		rctx.URLParams.Add("CARDID", id)
		rctx.URLParams.Add("BOARDID", "board_id")
	})

	rec := httptest.NewRecorder()

	getAPI(service).List(rec, req)

	res := rec.Result()

	service.AssertExpectations(t)

	expected := []*api.ListModel{&api.ListModel{
		ID:   string(model.ID),
		Name: model.Name,
	}}

	exp := strings.TrimSpace(string(toJson(t, expected)))
	act := strings.TrimSpace(string(body(t, res)))
	test.AssertExpAct(t, exp, act)
}

func getAPI(s board.Service) *api.API {
	return api.CreateAPI(s, &noop.Logger{})
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
