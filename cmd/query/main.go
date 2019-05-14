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

	l := shared.CreateLogger("[.QUERY.]")

	sess := shared.CreateSessionFactory()
	glob := shared.CreateSessionProvider(sess)
	prov := shared.CreateCopySessionProvider(glob)
	exec := shared.CreateOperationExecutor(prov)
	cfac := shared.CreateContextFactory(prov)
	rfac := shared.CreateRepositoryFactory(exec)
	sfac := shared.CreateServiceFactory(rfac)

	m := shared.ConfigureMux(cfac)

	exph := expvar.Handler()
	m.Get("/vars", func(w http.ResponseWriter, r *http.Request) { exph.ServeHTTP(w, r) })
	m.Get("/prof", pprof.Index)

	m.Route("/v1/api", func(r chi.Router) {
		boardRouter := chi.NewRouter()

		module := query.Module{
			Logger:         l,
			ServiceFactory: sfac,
			BoardRouter:    boardRouter,
			LaneRouter:     boardRouter,
			CardRouter:     boardRouter,
		}

		module.Boot()

		r.Mount("/board", boardRouter)
	})

	shared.StartMux(m, shared.GetPortOrDefault(8002), shared.CreateLogger("[..MUX..]"))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	glob.Provide().Close(false)

	l.Debugln("done")
}
