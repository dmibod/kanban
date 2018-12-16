package update

import (
	"log"
	"github.com/dmibod/kanban/tools/db/mongo"
	"github.com/dmibod/kanban/tools/mux"
)

// Boot - adds update module handlers to mux
func Boot(m mux.Mux){

	instance := func() interface{} {
		return &Card{}
	}

	repoFactory := mongo.New(mongo.WithDatabase("kanban"))

	m.Post("/post", &CreateCardHandler{ Repository: repoFactory.Create("cards", instance) })

	log.Println("Update module endpoints registered")
}