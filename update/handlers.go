package update

import (
	"log"
	"net/http"

	"github.com/dmibod/kanban/tools/db"
	"github.com/dmibod/kanban/tools/mux"
)

// CreateCard contains dependencies required by handler
type CreateCard struct {
	Repository db.Repository
}

// Parse parse request
func (h *CreateCard) Parse(r *http.Request) (interface{}, error) {
	card := &Card{}
	err := mux.JsonRequest(r, card)
	if err != nil {
		log.Println("Error parsing json", err)
	}
	return card, err
}

// Handle handles request
func (h *CreateCard) Handle(req interface{}) (interface{}, error) {
	card := req.(*Card)

	id, err := h.Repository.Create(card)
	if err != nil {
		log.Println("Error inserting document", err)
		return nil, err
	}

	res := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{id, true}

	return &res, nil
}
