package query

import (
	"log"
	"net/http"

	"github.com/dmibod/kanban/kernel"
)

// Card maps card to/from json at rest api level
type Card struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetCard contains dependencies required by handler
type GetCard struct {
	Service *CardService
}

// Parse parses Api request
func (h *GetCard) Parse(r *http.Request) (interface{}, error) {
	return r.FormValue("id"), nil
}

// Handle handles Api request
func (h *GetCard) Handle(req interface{}) (interface{}, error) {
	id := req.(string)

	log.Printf("GetCard request received: %v\n", id)

	card, err := h.Service.GetCardByID(kernel.Id(id))
	if err != nil {
		log.Println("Error getting card", err)
		return nil, err
	}

	return &Card{
		ID: string(card.ID),
		Name: card.Name,
	}, nil
}
