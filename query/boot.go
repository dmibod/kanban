package query

import (
	"net/http"
	"github.com/dmibod/kanban/tools/db/mongo"
	"github.com/dmibod/kanban/tools/mux"
)

func Boot(m mux.Mux){

	factory := func() interface{}{
		return &Card{}
	}
	env := &Env{ Db: mongo.New(mongo.WithDatabase("kanban"), mongo.WithCollection("cards"), mongo.WithFactory(factory)) }

	m.Handle("/get", http.HandlerFunc(env.GetCard))
}