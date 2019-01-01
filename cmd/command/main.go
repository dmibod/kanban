package main

import (
	"github.com/go-chi/chi"
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/command"
	"time"
)

func main() {
	c, cancel := context.WithCancel(context.Background())

	l := shared.CreateLogger("[COMMAND] ", true)
	m := shared.ConfigureMux()

	m.Route("/v1/api/command", func(r chi.Router) {
		router := chi.NewRouter()

		module := command.Module{Router: router, Logger: l}
		module.Boot()
	
		r.Mount("/", router)
	})

	shared.StartBus(c, shared.GetNameOrDefault("cmd"), shared.CreateLogger("[..BUS..] ", true))
	shared.StartMux(m, shared.GetPortOrDefault(8000), shared.CreateLogger("[..MUX..] ", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)

	shared.StopBus()
	
	l.Debugln("done")
}
