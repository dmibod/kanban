package main

import (
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/mux/http"
	"github.com/dmibod/kanban/update"
)

func main() {
	l := console.New(console.WithPrefix("[UPDATE.] "), console.WithDebug(true))
	m := http.New(http.WithPort(http.GetPortOrDefault(3003)))
	s := persistence.CreateService(l)
	f := mongo.CreateFactory(mongo.WithDatabase("kanban"), mongo.WithExecutor(s), mongo.WithLogger(l))

	update.Boot(m, f, l)

	m.Start()
}
