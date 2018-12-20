package main

import (
	"github.com/go-chi/chi"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	utils "github.com/dmibod/kanban/shared/tools/mux"
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

		module := query.Env{Logger: l, Factory: f, Mux: m }
		module.Boot()
	
		r.Mount("/", router)
	})

	utils.PrintRoutes(l, m)

	utils.StartMux(m, utils.GetPortOrDefault(3002))
}
