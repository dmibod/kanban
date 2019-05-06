package lane_test

import (
	"bytes"
	"context"
	"encoding/json"
	api "github.com/dmibod/kanban/query/lane"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services/lane"
	"github.com/dmibod/kanban/shared/services/lane/mocks"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"github.com/dmibod/kanban/shared/tools/test"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListByBoard(t *testing.T) {

	boardID := "board_id"
	laneID := "lane_id"

	model := &lane.ListModel{ID: kernel.ID(laneID), Name: "Sample"}

	service := &mocks.Service{}
	service.On("GetByBoardID", mock.Anything, kernel.ID(boardID)).Return([]*lane.ListModel{model}, nil).Once()

	req := toRequest(t, http.MethodGet, "", func(rctx *chi.Context) {
		rctx.URLParams.Add("BOARDID", boardID)
		rctx.URLParams.Add("LANEID", "")
	})

	rec := httptest.NewRecorder()

	getAPI(service).List(rec, req)

	res := rec.Result()

	service.AssertExpectations(t)

	expected := []*api.Lane{&api.Lane{
		ID:   string(model.ID),
		Name: model.Name,
	}}

	exp := strings.TrimSpace(string(toJson(t, expected)))
	act := strings.TrimSpace(string(body(t, res)))
	test.AssertExpAct(t, exp, act)
}

func TestListByParent(t *testing.T) {

	boardID := "board_id"
	parentID := "parent_id"
	laneID := "lane_id"

	model := &lane.ListModel{ID: kernel.ID(laneID), Name: "Sample"}

	service := &mocks.Service{}
	service.On("GetByLaneID", mock.Anything, kernel.ID(parentID).WithSet(kernel.ID(boardID))).Return([]*lane.ListModel{model}, nil).Once()

	req := toRequest(t, http.MethodGet, "", func(rctx *chi.Context) {
		rctx.URLParams.Add("BOARDID", boardID)
		rctx.URLParams.Add("LANEID", parentID)
	})

	rec := httptest.NewRecorder()

	getAPI(service).List(rec, req)

	res := rec.Result()

	service.AssertExpectations(t)

	expected := []*api.Lane{&api.Lane{
		ID:   string(model.ID),
		Name: model.Name,
	}}

	exp := strings.TrimSpace(string(toJson(t, expected)))
	act := strings.TrimSpace(string(body(t, res)))
	test.AssertExpAct(t, exp, act)
}

func TestOne(t *testing.T) {

	boardID := "board_id"
	laneID := "lane_id"

	model := &lane.Model{ID: kernel.ID(laneID), Name: "Sample"}

	service := &mocks.Service{}
	service.On("GetByID", mock.Anything, kernel.ID(laneID).WithSet(kernel.ID(boardID))).Return(model, nil).Once()

	req := toRequest(t, http.MethodGet, "", func(rctx *chi.Context) {
		rctx.URLParams.Add("BOARDID", boardID)
		rctx.URLParams.Add("LANEID", laneID)
	})

	rec := httptest.NewRecorder()

	getAPI(service).Get(rec, req)

	res := rec.Result()

	service.AssertExpectations(t)

	expected := &api.Lane{
		ID:   string(model.ID),
		Name: model.Name,
	}

	exp := strings.TrimSpace(string(toJson(t, expected)))
	act := strings.TrimSpace(string(body(t, res)))
	test.AssertExpAct(t, exp, act)
}

func getAPI(s lane.Service) *api.API {
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
