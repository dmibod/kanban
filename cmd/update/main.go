package main

import (
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/update"
	"github.com/go-chi/chi"
	"time"
)

func main() {

	l := shared.CreateLogger("[.UPDAT.]")

	sess := shared.CreateSessionFactory()
	glob := shared.CreateSessionProvider(sess)
	prov := shared.CreateCopySessionProvider(glob)
	exec := shared.CreateOperationExecutor(prov)
	cfac := shared.CreateContextFactory(prov)
	rfac := shared.CreateRepositoryFactory(exec)
	sfac := shared.CreateServiceFactory(rfac)

	m := shared.ConfigureMux(cfac)

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

	c, cancel := context.WithCancel(context.Background())
	shared.StartBus(c, shared.GetNameOrDefault("update"), shared.CreateLogger("[..BUS..] "))
	shared.StartMux(m, shared.GetPortOrDefault(8003), shared.CreateLogger("[..MUX..]"))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	glob.Provide().Close(false)

	cancel()

	time.Sleep(time.Second)

	shared.StopBus()

	l.Debugln("done")
}
