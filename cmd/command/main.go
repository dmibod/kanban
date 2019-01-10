package main

import (
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/command"
	"github.com/go-chi/chi"
	"time"
)

func main() {
	c, cancel := context.WithCancel(context.Background())

	l := shared.CreateLogger("[COMMAND] ")
	m := shared.ConfigureMux(nil)

	m.Route("/v1/api/command", func(r chi.Router) {
		router := chi.NewRouter()

		module := command.Module{Router: router, Logger: l}
		module.Boot()

		r.Mount("/", router)
	})

	shared.StartBus(c, shared.GetNameOrDefault("cmd"), shared.CreateLogger("[..BUS..] "))
	shared.StartMux(m, shared.GetPortOrDefault(8000), shared.CreateLogger("[..MUX..] "))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)

	shared.StopBus()

	l.Debugln("done")
}
