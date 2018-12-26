package main

import (
	"expvar"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"net/http"
	"net/http/pprof"
)

func main() {

	f := persistence.CreateFactory(
		persistence.CreateService(shared.CreateLogger("[BRK.MGO]", true)),
		shared.CreateLogger("[MONGO..]", true))

	m := shared.ConfigureMux()

	exph := expvar.Handler()
	m.Get("/vars", func(w http.ResponseWriter, r *http.Request) { exph.ServeHTTP(w, r) })
	m.Get("/prof", pprof.Index)

	module := query.Module{
		Logger:  shared.CreateLogger("[QUERY..]", true),
		Factory: services.CreateFactory(shared.CreateLogger("[SERVICE]", true), f),
		Mux:     m,
	}

	module.Boot(true)

	shared.StartMux(m, shared.GetPortOrDefault(8002), shared.CreateLogger("[..MUX..]", true))
}
