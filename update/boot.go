package update

import (
	"log"
	"github.com/dmibod/kanban/tools/db"
	"github.com/dmibod/kanban/tools/mux"
)

// Boot - adds update module handlers to mux
func Boot(m mux.Mux, f db.RepoFactory){

	instance := func() interface{} {
		return &Card{}
	}

	m.Post("/post", &CreateCardHandler{ Repository: f.Create("cards", instance) })

	log.Println("Update module endpoints registered")
}