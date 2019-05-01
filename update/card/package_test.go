package card_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"

	"github.com/dmibod/kanban/shared/tools/logger/noop"

	"github.com/dmibod/kanban/shared/services/card"

	"github.com/dmibod/kanban/shared/services/card/mocks"

	"github.com/dmibod/kanban/shared/kernel"

	updatecard "github.com/dmibod/kanban/update/card"
)

func TestCardAPI(t *testing.T) {
	id := "5c16dd24c7ee6e5dcf626266"

	payload := &updatecard.Card{ID: id, Name: "Sample"}

	param := &card.CreateModel{Name: "Sample"}
	model := &card.Model{ID: kernel.ID(id), Name: "Sample"}

	service := &mocks.Service{}
	service.On("Create", mock.Anything, kernel.ID("board_id"), param).Return(kernel.ID(id), nil).Once()
	service.On("GetByID", mock.Anything, kernel.ID(id).WithSet("board_id")).Return(model, nil).Once()

	req := toJsonRequest(t, http.MethodPost, "http://localhost/v1/api/board_id/cards/", payload, func(rctx *chi.Context) {
		rctx.URLParams.Add("BOARDID", "board_id")
	})
	res := httptest.NewRecorder()

	getAPI(service).CreateCard(res, req)

	service.AssertExpectations(t)

	expected := &updatecard.Card{
		ID:   id,
		Name: payload.Name,
	}

	exp := strings.TrimSpace(string(toJson(t, &expected)))
	act := strings.TrimSpace(res.Body.String())
	test.AssertExpAct(t, exp, act)
}

func getAPI(s card.Service) *updatecard.API {
	return updatecard.CreateAPI(s, &noop.Logger{})
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
