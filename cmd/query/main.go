package main

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/query"
)

func main() {

	l := shared.CreateLogger("[.QUERY.]", true)
	m := shared.ConfigureMux()

	exph := expvar.Handler()
	m.Get("/vars", func(w http.ResponseWriter, r *http.Request) { exph.ServeHTTP(w, r) })
	m.Get("/prof", pprof.Index)

	m.Route("/v1/api", func(r chi.Router) {
		cardRouter := chi.NewRouter()

		module := query.Module{
			Logger:     l,
			Factory:    shared.CreateServiceFactory(),
			CardRouter: cardRouter,
		}

		module.Boot()

		r.Mount("/card", cardRouter)
	})

	shared.StartMux(m, shared.GetPortOrDefault(8002), shared.CreateLogger("[..MUX..]", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")
	l.Debugln("done")
}
