package query

import (
	"log"

	"github.com/dmibod/kanban/tools/db/mongo"
	"github.com/dmibod/kanban/tools/mux"
)

// Boot installs Query module http handlers to mux
func Boot(m mux.Mux) {

	repoFactory := mongo.New(mongo.WithDatabase("kanban"))

	m.Get("/get", &GetCardHandler{Service: CreateCardService(CreateCardRepository(repoFactory))})

	log.Println("Query module endpoints registered")
}
