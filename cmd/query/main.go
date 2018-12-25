package main

import (
	"expvar"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"net/http"
	"net/http/pprof"
)

func main() {

	l := console.New(
		console.WithPrefix("[QUERY..] "),
		console.WithDebug(true))

	f := mongo.CreateFactory(
		"kanban",
		persistence.CreateService(l),
		l)

	m := mux.ConfigureMux()

	exph := expvar.Handler()
	m.Get("/vars", func(w http.ResponseWriter, r *http.Request) { exph.ServeHTTP(w, r) })
	m.Get("/prof", pprof.Index)

	module := query.Module{Logger: l, Factory: services.CreateFactory(l, f), Mux: m}
	module.Boot(true)

	mux.StartMux(m, mux.GetPortOrDefault(8002), l)
}
