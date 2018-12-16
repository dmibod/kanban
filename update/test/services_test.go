package test

import (
	"github.com/dmibod/kanban/tools/mux"
	"encoding/json"
	"bytes"
	"github.com/dmibod/kanban/kernel"
	"github.com/dmibod/kanban/update"
	"net/http/httptest"
	"net/http"
	"testing"
	mock "github.com/dmibod/kanban/tools/db/mocks"
)

func TestCreateCard(t *testing.T) {

	id := "000"
	card := update.Card{ Id: kernel.Id(id), Name: "Sample" }

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
	 
	h := mux.Handle(&update.CreateCard{ Repository: repo })
	h.ServeHTTP(w, r)

	repo.AssertExpectations(t)
}
