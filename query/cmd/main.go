package main

import (
	"net/http"
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/tools/db/mongo"
)

func main() {
	factory := func() interface{}{
		return &query.Card{}
	}
	env := &query.Env{ Db: mongo.New(mongo.WithDatabase("kanban"), mongo.WithCollection("cards"), mongo.WithFactory(factory)) }

	http.HandleFunc("/card", env.GetCard)
	
	http.ListenAndServe(":3002", nil)
}