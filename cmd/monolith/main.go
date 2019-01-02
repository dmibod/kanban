package main

import (
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"context"
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
	l := shared.CreateLogger("[KANBAN.] ", true)
	c, cancel := context.WithCancel(context.Background())

	sess := shared.CreateSessionFactory()

	bootWorkers(c, sess)
	bootWeb(sess)

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)

	shared.StopBus()

	l.Debugln("done")
}

func bootWorkers(ctx context.Context, sess mongo.SessionFactory) {
	prov := shared.CreateSessionProvider(sess)
	exec := shared.CreateExecutor(prov)
	cfac := shared.CreateContextFactory(prov)
	rfac := shared.CreateRepositoryFactory(exec)
	sfac := shared.CreateServiceFactory(rfac)

	ctx, err := cfac.Context(ctx)
	if err != nil {
		panic(err)
	}

	boot(&process.Module{Context: ctx, ServiceFactory: sfac})

	shared.StartBus(ctx, shared.GetNameOrDefault("mono"), shared.CreateLogger("[..BUS..] ", true))
}

func bootWeb(sess mongo.SessionFactory) {
	prov := shared.CreateSessionProvider(sess)
	cfac := shared.CreateContextFactory(prov)

	m := shared.ConfigureMux(cfac)

	exec := shared.CreateExecutor(prov)
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

	shared.StartMux(m, shared.GetPortOrDefault(3000), shared.CreateLogger("[..MUX..] ", true))
}

func boot(b interface{ Boot() }) {
	b.Boot()
}
