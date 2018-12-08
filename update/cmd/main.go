package main

import (
	"net/http"
	"github.com/dmibod/kanban/update"
)

func main() {
	env := &update.Env{ Db: "todo"/*mongo.New()*/ }

	http.HandleFunc("/", env.CreateCard)
	http.ListenAndServe(":3003", nil)
}