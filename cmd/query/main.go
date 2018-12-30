package main

import (
	"github.com/go-chi/chi"
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

	m.Route("/v1/api", func(r chi.Router) {
		card := chi.NewRouter()

		module := query.Module{
			Logger:  shared.CreateLogger("[.QUERY.]", true),
			Factory: shared.CreateServiceFactory(),
			Card:    card,
		}

		module.Boot()

		r.Mount("/card", card)
	})


	shared.StartMux(m, shared.GetPortOrDefault(8002), shared.CreateLogger("[..MUX..]", true))
}
