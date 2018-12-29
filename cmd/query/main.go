package main

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/query"
)

func main() {

	m := shared.ConfigureMux()

	exph := expvar.Handler()
	m.Get("/vars", func(w http.ResponseWriter, r *http.Request) { exph.ServeHTTP(w, r) })
	m.Get("/prof", pprof.Index)

	module := query.Module{
		Logger:  shared.CreateLogger("[QUERY..]", true),
		Factory: shared.CreateServiceFactory(),
		Mux:     m,
	}

	module.Boot(true)

	shared.StartMux(m, shared.GetPortOrDefault(8002), shared.CreateLogger("[..MUX..]", true))
}
