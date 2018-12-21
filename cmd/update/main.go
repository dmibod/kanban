package main

import (
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	utils "github.com/dmibod/kanban/shared/tools/mux"
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

	m := utils.ConfigureMux()

	module := update.Module{Logger: l, Factory: services.CreateFactory(l, f), Mux: m}
	module.Standalone()

	utils.PrintRoutes(l, m)
	utils.StartMux(m, utils.GetPortOrDefault(3003))
}
