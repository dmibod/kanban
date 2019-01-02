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
	l := shared.CreateLogger("[KANBAN.] ", true)

	slog := shared.CreateLogger("[SESSION] ", true)
	sess := mongo.CreateSessionFactory(mongo.WithLogger(slog))
	exec := shared.CreateExecutor(sess)
	rfac := shared.CreateRepositoryFactory(exec)
	sfac := shared.CreateServiceFactory(rfac)

	cfac := mongo.CreateContextFactory(sess, slog)
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

	c, cancel := context.WithCancel(context.Background())

	boot(&process.Module{Context: c, ServiceFactory: sfac, ContextFactory: cfac})

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
