package query

import (
	"log"
	"net/http"

	"github.com/dmibod/kanban/kernel"
	"github.com/dmibod/kanban/tools/mux"
)

// Card maps card to/from json at rest api level
type Card struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetCardHandler contains dependencies required handler
type GetCardHandler struct {
	Service *CardService
}

// GetCard implements /get?id= method
func (h *GetCardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	log.Printf("GetCard request received: %v\n", id)

	card, err := h.Service.GetCardByID(kernel.Id(id))
	if err != nil {
		mux.ErrorResponse(w, http.StatusInternalServerError)
		log.Println("Error getting card", err)
		return
	}

	mux.JsonResponse(w, &Card{
		ID: string(card.ID),
		Name: card.Name,
	})
}
