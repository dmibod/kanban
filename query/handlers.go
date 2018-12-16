package query

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dmibod/kanban/kernel"
)

// Env contains dependencies required by http handlers
type GetCardHandler struct {
	Service *CardService
}

// GetCard implements /get?id= method
func (h *GetCardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	log.Printf("GetCard request received: %v\n", id)

	card, err := h.Service.GetCardByID(kernel.Id(id))
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println("Error getting card", err)
		return
	}

	enc := json.NewEncoder(w)
	enc.Encode(Card{
		ID: string(card.ID),
		Name: card.Name,
	})
}
