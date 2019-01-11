package update_test

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

	"github.com/dmibod/kanban/shared/services"

	"github.com/dmibod/kanban/shared/services/mocks"

	"github.com/dmibod/kanban/shared/kernel"

	"github.com/dmibod/kanban/update"
)

func TestCardAPI(t *testing.T) {
	id := "5c16dd24c7ee6e5dcf626266"

	payload := &update.Card{ID: id, Name: "Sample"}

	param := &services.CardPayload{Name: "Sample"}
	model := &services.CardModel{ID: kernel.ID(id), Name: "Sample"}

	service := &mocks.CardService{}
	service.On("Create", mock.Anything, param).Return(model, nil).Once()

	req := toJsonRequest(t, http.MethodPost, "http://localhost/v1/api/card/", payload)
	res := httptest.NewRecorder()

	getAPI(service).CreateCard(res, req)

	service.AssertExpectations(t)

	expected := &update.Card{
		ID:   id,
		Name: payload.Name,
	}

	exp := strings.TrimSpace(string(toJson(t, &expected)))
	act := strings.TrimSpace(res.Body.String())
	test.AssertExpAct(t, exp, act)
}

func getAPI(s services.CardService) *update.CardAPI {
	return update.CreateCardAPI(s, &noop.Logger{})
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
