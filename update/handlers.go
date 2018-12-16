package update

import (
	"log"
	"net/http"

	"github.com/dmibod/kanban/tools/db"
	"github.com/dmibod/kanban/tools/mux"
)

type CreateCardHandler struct {
	Repository db.Repository
}

func (h *CreateCardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	card := &Card{}

	jsonErr := mux.JsonRequest(r, card)
	if jsonErr != nil {
		mux.ErrorResponse(w, http.StatusInternalServerError)
		log.Println("Error parsing json", jsonErr)
		return
	}

	id, dbErr := h.Repository.Create(card)
	if dbErr != nil {
		mux.ErrorResponse(w, http.StatusInternalServerError)
		log.Println("Error inserting document", dbErr)
		return
	}

	d := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{id, true}

	mux.JsonResponse(w, &d)
}
