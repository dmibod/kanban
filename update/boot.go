package update

import (
	"log"
	"net/http"
	"github.com/dmibod/kanban/tools/db/mongo"
	"github.com/dmibod/kanban/tools/mux"
)

// Boot - adds update module handlers to mux
func Boot(m mux.Mux){

	instance := func() interface{} {
		return &Card{}
	}

	repoFactory := mongo.New(mongo.WithDatabase("kanban"))

	env := &Env{ Repository: repoFactory.Create("cards", instance) }

	http.HandleFunc("/post", env.CreateCard)

	log.Println("Update module endpoints registered")
}