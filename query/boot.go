package query

import (
	"log"

	"github.com/dmibod/kanban/tools/db"
	"github.com/dmibod/kanban/tools/mux"
)

// Boot installs Query module http handlers to mux
func Boot(m mux.Mux, f db.RepoFactory) {

	m.Get("/get", &GetCardHandler{Service: CreateCardService(CreateCardRepository(f))})

	log.Println("Query module endpoints registered")
}
