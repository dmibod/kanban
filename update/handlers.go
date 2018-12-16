package update

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dmibod/kanban/tools/db"
)

type CreateCardHandler struct {
	Repository db.Repository
}

func (h *CreateCardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, readErr := ioutil.ReadAll(r.Body)

	if readErr != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println("Error reading body", readErr)
		return
	}

	card := Card{}
	jsonErr := json.Unmarshal(body, &card)

	if jsonErr != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println("Error parsing json", jsonErr)
		return
	}

	id, dbErr := h.Repository.Create(&card)

	if dbErr != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println("Error inserting document", dbErr)
		return
	}

	enc := json.NewEncoder(w)
	d := struct {
		Id      string `json:"id"`
		Success bool   `json:"success"`
	}{id, true}

	enc.Encode(d)
}
