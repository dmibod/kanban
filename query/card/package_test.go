package card_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dmibod/kanban/shared/services/card/mocks"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	api "github.com/dmibod/kanban/query/card"

	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/stretchr/testify/mock"

	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services/card"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

func TestList(t *testing.T) {

	boardID := "board_id"
	laneID := "lane_id"
	cardID := "card_id"

	model := &card.Model{ID: kernel.ID(cardID), Name: "Sample"}

	service := &mocks.Service{}
	service.On("GetByLaneID", mock.Anything, kernel.ID(laneID).WithSet(kernel.ID(boardID))).Return([]*card.Model{model}, nil).Once()

	req := toRequest(t, http.MethodGet, "", func(rctx *chi.Context) {
		rctx.URLParams.Add("LANEID", laneID)
		rctx.URLParams.Add("BOARDID", boardID)
	})

	rec := httptest.NewRecorder()

	getAPI(service).List(rec, req)

	res := rec.Result()

	service.AssertExpectations(t)

	expected := []*api.Card{&api.Card{
		ID:   string(model.ID),
		Name: model.Name,
	}}

	exp := strings.TrimSpace(string(toJson(t, expected)))
	act := strings.TrimSpace(string(body(t, res)))
	test.AssertExpAct(t, exp, act)
}

func TestOne(t *testing.T) {

	boardID := "board_id"
	cardID := "card_id"

	model := &card.Model{ID: kernel.ID(cardID), Name: "Sample"}

	service := &mocks.Service{}
	service.On("GetByID", mock.Anything, kernel.ID(cardID).WithSet(kernel.ID(boardID))).Return(model, nil).Once()

	req := toRequest(t, http.MethodGet, "", func(rctx *chi.Context) {
		rctx.URLParams.Add("CARDID", cardID)
		rctx.URLParams.Add("BOARDID", boardID)
	})

	rec := httptest.NewRecorder()

	getAPI(service).Get(rec, req)

	res := rec.Result()

	service.AssertExpectations(t)

	expected := &api.Card{
		ID:   string(model.ID),
		Name: model.Name,
	}

	exp := strings.TrimSpace(string(toJson(t, expected)))
	act := strings.TrimSpace(string(body(t, res)))
	test.AssertExpAct(t, exp, act)
}

func getAPI(s card.Service) *api.API {
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
