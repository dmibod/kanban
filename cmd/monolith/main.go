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
	c, cancel := context.WithCancel(context.Background())

	boot(&process.Module{Logger: shared.CreateLogger("[PROCESS] ", true), Context: c})

	l := shared.CreateLogger("[KANBAN] ", true)
	m := shared.ConfigureMux()
	s := shared.CreateServiceFactory()

	boot(&notify.Module{Logger: shared.CreateLogger("[.NOTIF.] ", true), Mux: m})

	m.Route("/v1/api", func(r chi.Router) {
		commandRouter := chi.NewRouter()
		boardRouter := chi.NewRouter()
		cardRouter := chi.NewRouter()

		boot(&command.Module{Router: commandRouter})
		boot(&query.Module{BoardRouter: boardRouter, CardRouter: cardRouter, Factory: s})
		boot(&update.Module{BoardRouter: boardRouter, CardRouter: cardRouter, Factory: s})

		r.Mount("/command", commandRouter)
		r.Mount("/board", boardRouter)
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
