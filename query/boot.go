package query

import (
	"log"
	"net/http"
	"github.com/dmibod/kanban/tools/db/mongo"
	"github.com/dmibod/kanban/tools/mux"
)

func Boot(m mux.Mux){

	instance := func() interface{}{
		return &Card{}
	}

	repoFactory := mongo.New(mongo.WithDatabase("kanban"))

	env := &Env{ Repository: repoFactory.Create("cards", instance) }

	m.Handle("/get", http.HandlerFunc(env.GetCard))

	log.Println("Query module endpoints registered")
}