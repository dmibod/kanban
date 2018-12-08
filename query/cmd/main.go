package main

import (
	"net/http"
	"github.com/dmibod/kanban/query"
)

func main() {
	env := &query.Env{ Db: "todo"/*mongo.New()*/ }

	http.HandleFunc("/", env.GetCards)
	http.ListenAndServe(":3002", nil)
}