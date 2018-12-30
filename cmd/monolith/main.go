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

	boot(&command.Module{Logger: shared.CreateLogger("[COMMAND] ", true), Mux: m})
	boot(&notify.Module{Logger: shared.CreateLogger("[.NOTIF.] ", true), Mux: m})

	m.Route("/v1/api", func(r chi.Router) {
		board := chi.NewRouter()
		card := chi.NewRouter()

		boot(&query.Module{Logger: shared.CreateLogger("[.QUERY.] ", true), Card: card, Factory: s})
		boot(&update.Module{Logger: shared.CreateLogger("[.UPDAT.] ", true), Board: board, Card: card, Factory: s})

		r.Mount("/board", board)
		r.Mount("/card", card)
	})

	shared.StartBus(c, shared.GetNameOrDefault("mono"), shared.CreateLogger("[..BUS..] ", true))
	shared.StartMux(m, shared.GetPortOrDefault(8000), shared.CreateLogger("[..MUX..] ", true))

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
