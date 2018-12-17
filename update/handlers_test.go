package update_test

import (
	"strings"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmibod/kanban/shared/kernel"

	logm "github.com/dmibod/kanban/shared/tools/log/mocks"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/dmibod/kanban/update"
)

type mockCardService struct {
	Id kernel.Id
}

func (s *mockCardService) CreateCard(p *update.CardPayload) (kernel.Id, error) {
	return s.Id, nil
}

func TestCreateCard(t *testing.T) {

	id := "000"
	card := update.Card{ID: id, Name: "Sample"}

	payload, jsonErr := json.Marshal(&card)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	req, reqErr := http.NewRequest("POST", "http://localhost/post", bytes.NewBuffer(payload))
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	w := httptest.NewRecorder()

	mux.Handle(update.CreateCreateCardHandler(&logm.Logger{}, &mockCardService{kernel.Id(id)})).ServeHTTP(w, req)

	got := strings.TrimSpace(string(w.Body.Bytes()))

	expected := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{id, true}

	resp, respErr := json.Marshal(&expected)
	if respErr != nil {
		t.Fatal(respErr)
	}
	
	want := strings.TrimSpace(string(resp))

	if got != want {
		t.Fatalf("Wrong response\nwant: %v\ngot: %v", string(want), string(got))
	}
}
