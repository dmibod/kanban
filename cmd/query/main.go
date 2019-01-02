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

	e, p := shared.CreateDatabaseServices()
	rfac := shared.CreateRepositoryFactory(e)
	sfac := shared.CreateServiceFactory(rfac)

	m := shared.ConfigureMux(p)

	exph := expvar.Handler()
	m.Get("/vars", func(w http.ResponseWriter, r *http.Request) { exph.ServeHTTP(w, r) })
	m.Get("/prof", pprof.Index)

	m.Route("/v1/api", func(r chi.Router) {
		boardRouter := chi.NewRouter()
		laneRouter := chi.NewRouter()
		cardRouter := chi.NewRouter()

		module := query.Module{
			Logger:      l,
			Factory:     sfac,
			BoardRouter: boardRouter,
			LaneRouter:  laneRouter,
			CardRouter:  cardRouter,
		}

		module.Boot()

		r.Mount("/board", boardRouter)
		r.Mount("/lane", laneRouter)
		r.Mount("/card", cardRouter)
	})

	shared.StartMux(m, shared.GetPortOrDefault(8002), shared.CreateLogger("[..MUX..]", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")
	l.Debugln("done")
}
