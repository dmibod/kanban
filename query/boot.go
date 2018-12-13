package query

import (
	"log"
	"net/http"

	"github.com/dmibod/kanban/tools/db/mongo"
	"github.com/dmibod/kanban/tools/mux"
)

// Boot installs Query module http handlers to mux
func Boot(m mux.Mux) {

	repoFactory := mongo.New(mongo.WithDatabase("kanban"))

	env := &Env{Service: CreateCardService(CreateCardRepository(repoFactory))}

	m.Handle("/get", http.HandlerFunc(env.GetCard))

	log.Println("Query module endpoints registered")
}
