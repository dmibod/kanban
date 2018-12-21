package main

import (
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/update"
)

func main() {
	l := console.New(
		console.WithPrefix("[UPDATE.] "),
		console.WithDebug(true))

	f := mongo.CreateFactory(
		mongo.WithDatabase("kanban"),
		mongo.WithExecutor(persistence.CreateService(l)),
		mongo.WithLogger(l))

	m := mux.ConfigureMux()

	module := update.Module{Logger: l, Factory: services.CreateFactory(l, f), Mux: m}
	module.Boot(true)

	mux.PrintRoutes(l, m)
	mux.StartMux(m, mux.GetPortOrDefault(3003))
}
