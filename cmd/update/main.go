package main

import (
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/update"
	"github.com/go-chi/chi"
)

func main() {

	l := shared.CreateLogger("[.UPDAT.]", true)

	slog := shared.CreateLogger("[SESSION] ", true)
	sess := mongo.CreateSessionFactory(mongo.WithLogger(slog))
	exec := shared.CreateExecutor(sess)
	rfac := shared.CreateRepositoryFactory(exec)
	sfac := shared.CreateServiceFactory(rfac)

	m := shared.ConfigureMux(mongo.CreateContextFactory(sess, slog))

	m.Route("/v1/api", func(r chi.Router) {
		boardRouter := chi.NewRouter()
		laneRouter := chi.NewRouter()
		cardRouter := chi.NewRouter()

		module := update.Module{
			Logger:         l,
			ServiceFactory: sfac,
			BoardRouter:    boardRouter,
			LaneRouter:     laneRouter,
			CardRouter:     cardRouter,
		}

		module.Boot()

		r.Mount("/board", boardRouter)
		r.Mount("/lane", laneRouter)
		r.Mount("/card", cardRouter)
	})

	shared.StartMux(m, shared.GetPortOrDefault(8003), shared.CreateLogger("[..MUX..]", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")
	l.Debugln("done")
}
