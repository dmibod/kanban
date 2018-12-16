package command

import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/dmibod/kanban/tools/mux"
)

// PostCommandsHandler handles PostCommands end point
type PostCommandsHandler struct {
	CommandQueue chan<- []byte
}

func (h *PostCommandsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	commands := []Command{}

	jsonErr := mux.JsonRequest(r, &commands)
	if jsonErr != nil {
		mux.ErrorResponse(w, http.StatusInternalServerError)
		log.Println("Error parsing json", jsonErr)
		return
	}

	log.Printf("Commands received: %+v\n", commands);

	m, msgErr := json.Marshal(commands)
	if msgErr != nil {
		mux.ErrorResponse(w, http.StatusInternalServerError)
		log.Println("Error marshalling commands", msgErr)
		return
	}

	h.CommandQueue <- m

	d := struct {
		Count int `json:"count"`
		Success bool `json:"success"`
	}{len(commands),true}

	mux.JsonResponse(w, &d)

	log.Printf("Commands sent: %+v\n", len(commands))
}
