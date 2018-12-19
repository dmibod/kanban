package main

import (
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/mux/http"
)

func main() {
	m := http.New(http.WithPort(http.GetPortOrDefault(3002)))
	s := persistence.CreateDatabaseService(nil)
	f := mongo.CreateFactory(mongo.WithDatabase("kanban"), mongo.WithExecutor(s))

	query.Boot(m, f)

	m.Start()
}
