package main

import (
	"net/http"
	"github.com/dmibod/kanban/update"
	"github.com/dmibod/kanban/tools/db/mongo"
)

func main() {
	factory := func() interface{}{
		return &update.Card{}
	}
	env := &update.Env{ Db: mongo.New(mongo.WithDatabase("kanban"), mongo.WithCollection("cards"), mongo.WithFactory(factory)) }

	http.HandleFunc("/", env.CreateCard)
	http.ListenAndServe(":3003", nil)
}