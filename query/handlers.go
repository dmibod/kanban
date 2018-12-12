package query

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/dmibod/kanban/tools/db"
)

type Env struct {
	Repository db.Repository
}

func (e *Env) GetCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")

	log.Printf("GetCard request received: %v\n", id)

	card, err := e.Repository.FindById(id)
	if (err != nil){
		http.Error(w, http.StatusText(500), 500)
		log.Println("Error getting card", err)
		return
	}

	enc := json.NewEncoder(w)
	enc.Encode(card)
}
