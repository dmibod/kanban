package update_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmibod/kanban/shared/kernel"

	logm "github.com/dmibod/kanban/shared/tools/log/mocks"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/dmibod/kanban/update"
)

type cardServiceMock struct {
	Id kernel.Id
}

func (s *cardServiceMock) CreateCard(p *update.CardPayload) (kernel.Id, error) {
	return s.Id, nil
}

func TestCreateCard(t *testing.T) {

	payload := update.Card{ID: "000", Name: "Sample"}

	req := toJsonRequest(t, http.MethodPost, "http://localhost/post", &payload)
	res := httptest.NewRecorder()

	service := &cardServiceMock{kernel.Id(payload.ID)}
	handler := update.CreateCreateCardHandler(&logm.Logger{}, service)

	mux.Handle(handler).ServeHTTP(res, req)

	expected := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{payload.ID, true}

	exp := strings.TrimSpace(string(toJson(t, &expected)))
	act := strings.TrimSpace(res.Body.String())

	assertf(t, act == exp, "Wrong response\nwant: %v\ngot: %v", exp, act)
}
