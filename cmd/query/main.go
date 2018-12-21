package main

import (
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	utils "github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/chi"
)

func main() {

	l := console.New(
		console.WithPrefix("[QUERY..] "),
		console.WithDebug(true))

	f := mongo.CreateFactory(
		mongo.WithDatabase("kanban"),
		mongo.WithExecutor(persistence.CreateService(l)),
		mongo.WithLogger(l))

	m := utils.ConfigureMux()

	m.Route("/v1/api/card", func(r chi.Router) {
		router := chi.NewRouter()

		module := query.Module{Logger: l, Factory: services.CreateFactory(l, f), Mux: m}
		module.Boot()

		r.Mount("/", router)
	})

	utils.PrintRoutes(l, m)

	utils.StartMux(m, utils.GetPortOrDefault(3002))
}
