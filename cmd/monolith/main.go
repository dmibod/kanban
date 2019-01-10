package main

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"time"

	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/notify"
	"github.com/dmibod/kanban/process"
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/update"
)

func main() {
	l := shared.CreateLogger("[KANBAN.] ")
	c, cancel := context.WithCancel(context.Background())

	sess := shared.CreateSessionFactory()
	prov := shared.CreateSessionProvider(sess)

	bootWks(c, prov)
	bootWeb(prov)

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	prov.Provide().Close(false)

	cancel()

	time.Sleep(time.Second)

	shared.StopBus()

	l.Debugln("done")
}

func bootWks(ctx context.Context, glob mongo.SessionProvider) {
	prov := shared.CreateCopySessionProvider(glob)
	exec := shared.CreateOperationExecutor(prov)
	rfac := shared.CreateRepositoryFactory(exec)

	cfac := shared.CreateContextFactory(prov)
	sfac := shared.CreateServiceFactory(rfac)

	boot(&process.Module{ContextFactory: cfac, ServiceFactory: sfac})

	shared.StartBus(ctx, shared.GetNameOrDefault("mono"), shared.CreateLogger("[..BUS..] "))
}

func bootWeb(glob mongo.SessionProvider) {
	prov := shared.CreateCopySessionProvider(glob)
	cfac := shared.CreateContextFactory(prov)

	m := shared.ConfigureMux(cfac)

	exec := shared.CreateOperationExecutor(prov)
	rfac := shared.CreateRepositoryFactory(exec)
	sfac := shared.CreateServiceFactory(rfac)

	m.Route("/v1/api", func(r chi.Router) {
		commandRouter := chi.NewRouter()
		notifyRouter := chi.NewRouter()
		boardRouter := chi.NewRouter()
		laneRouter := chi.NewRouter()
		cardRouter := chi.NewRouter()

		boot(&command.Module{Router: commandRouter})
		boot(&notify.Module{Router: notifyRouter})
		boot(&query.Module{BoardRouter: boardRouter, LaneRouter: laneRouter, CardRouter: cardRouter, ServiceFactory: sfac})
		boot(&update.Module{BoardRouter: boardRouter, LaneRouter: laneRouter, CardRouter: cardRouter, ServiceFactory: sfac})

		r.Mount("/command", commandRouter)
		r.Mount("/notify", notifyRouter)
		r.Mount("/board", boardRouter)
		r.Mount("/lane", laneRouter)
		r.Mount("/card", cardRouter)
	})

	shared.StartMux(m, shared.GetPortOrDefault(3001), shared.CreateLogger("[..MUX..] "))
}

func boot(b interface{ Boot() }) {
	b.Boot()
}
