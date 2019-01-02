package main

import (
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

	sess := shared.CreateSessionFactory()
	prov := shared.CreateSessionProvider(sess)
	exec := shared.CreateExecutor(prov)
	cfac := shared.CreateContextFactory(prov)
	rfac := shared.CreateRepositoryFactory(exec)
	sfac := shared.CreateServiceFactory(rfac)

	c, cancel := context.WithCancel(context.Background())

	ctx, err := cfac.Context(c)
	if err != nil {
		l.Errorln(err)
		cancel()
		return
	}

	boot(&process.Module{Context: ctx, ServiceFactory: sfac})

	prov = shared.CreateSessionProvider(sess)
	exec = shared.CreateExecutor(prov)
	cfac = shared.CreateContextFactory(prov)
	rfac = shared.CreateRepositoryFactory(exec)
	sfac = shared.CreateServiceFactory(rfac)

	m := shared.ConfigureMux(cfac)

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

	shared.StartBus(c, shared.GetNameOrDefault("mono"), shared.CreateLogger("[..BUS..] ", true))
	shared.StartMux(m, shared.GetPortOrDefault(3000), shared.CreateLogger("[..MUX..] ", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)

	shared.StopBus()

	l.Debugln("done")
}

func boot(b interface{ Boot() }) {
	b.Boot()
}
