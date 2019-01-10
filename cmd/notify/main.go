package main

import (
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/notify"
	"github.com/go-chi/chi"
	"time"
)

func main() {
	c, cancel := context.WithCancel(context.Background())

	l := shared.CreateLogger("[.NOTIF.] ")
	m := shared.ConfigureMux(nil)

	m.Route("/v1/api", func(r chi.Router) {
		router := chi.NewRouter()

		module := notify.Module{Router: router, Logger: l}
		module.Boot()

		r.Mount("/notify", router)
	})

	shared.StartBus(c, shared.GetNameOrDefault("notify"), shared.CreateLogger("[..BUS..] "))
	shared.StartMux(m, shared.GetPortOrDefault(8001), shared.CreateLogger("[..MUX..] "))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)

	shared.StopBus()

	l.Debugln("done")
}
