package main

import (
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/cmd/shared"
)

func main() {

	l := console.New(
		console.WithPrefix("[QUERY..] "),
		console.WithDebug(true))

	f := mongo.CreateFactory(
		mongo.WithDatabase("kanban"),
		mongo.WithExecutor(persistence.CreateService(l)),
		mongo.WithLogger(l))

	m := mux.ConfigureMux()

	module := query.Module{Logger: l, Factory: services.CreateFactory(l, f), Mux: m}
	module.Boot(true)

	mux.PrintRoutes(l, m)
	mux.StartMux(m, mux.GetPortOrDefault(3002))
}
