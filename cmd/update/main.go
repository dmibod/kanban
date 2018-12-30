package main

import (
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/update"
	"github.com/go-chi/chi"
)

func main() {

	l := shared.CreateLogger("[.UPDAT.]", true)
	m := shared.ConfigureMux()

	m.Route("/v1/api", func(r chi.Router) {
		board := chi.NewRouter()
		card := chi.NewRouter()

		module := update.Module{
			Logger:  l,
			Factory: shared.CreateServiceFactory(),
			Board:   board,
			Card:    card,
		}

		module.Boot()

		r.Mount("/board", board)
		r.Mount("/card", card)
	})

	shared.StartMux(m, shared.GetPortOrDefault(8003), shared.CreateLogger("[..MUX..]", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")
	l.Debugln("done")
}
