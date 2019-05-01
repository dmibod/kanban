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

	boardID := "board_id"

	model := &board.ListModel{ID: kernel.ID(boardID), Name: "Sample"}

	service := &mocks.Service{}
	service.On("GetByOwner", mock.Anything, mock.Anything).Return([]*board.ListModel{model}, nil).Once()

	req := toRequest(t, http.MethodGet, "")

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

func TestOne(t *testing.T) {

	boardID := "board_id"

	model := &board.Model{ID: kernel.ID(boardID), Name: "Sample"}

	service := &mocks.Service{}
	service.On("GetByID", mock.Anything, kernel.ID(boardID)).Return(model, nil).Once()

	req := toRequest(t, http.MethodGet, "", func(rctx *chi.Context) {
		rctx.URLParams.Add("BOARDID", boardID)
	})

	rec := httptest.NewRecorder()

	getAPI(service).Get(rec, req)

	res := rec.Result()

	service.AssertExpectations(t)

	expected := &api.Model{
		ID:   string(model.ID),
		Name: model.Name,
	}

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
