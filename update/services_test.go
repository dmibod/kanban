package update_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmibod/kanban/shared/kernel"
	mock "github.com/dmibod/kanban/shared/tools/db/mocks"
	"github.com/dmibod/kanban/shared/tools/log/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/dmibod/kanban/update"
)

func TestCreateCard(t *testing.T) {

	id := "000"
	card := update.Card{Id: kernel.Id(id), Name: "Sample"}

	jsonPayload, jsonErr := json.Marshal(&card)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	r, err := http.NewRequest("POST", "http://localhost/post", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	repo := &mock.Repository{}
	repo.On("Create", &card).Return(id, nil).Once()

	h := mux.Handle(&update.CreateCard{Logger: logger.New(), Repository: repo})
	h.ServeHTTP(w, r)

	repo.AssertExpectations(t)
}
