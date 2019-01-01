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
		boardRouter := chi.NewRouter()
		cardRouter := chi.NewRouter()

		module := update.Module{
			Logger:      l,
			Factory:     shared.CreateServiceFactory(),
			BoardRouter: boardRouter,
			CardRouter:  cardRouter,
		}

		module.Boot()

		r.Mount("/board", boardRouter)
		r.Mount("/card", cardRouter)
	})

	shared.StartMux(m, shared.GetPortOrDefault(8003), shared.CreateLogger("[..MUX..]", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")
	l.Debugln("done")
}
