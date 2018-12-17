package main

import (
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/shared/tools/mux/http"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

func main() {
	m := http.New(http.WithPort(http.GetPortOrDefault(3002)))
	f := mongo.New(mongo.WithDatabase("kanban"))

	query.Boot(m, f)

	m.Start()
}
