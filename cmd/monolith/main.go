package main

import (
	"context"
	"github.com/dmibod/kanban/shared/services"
	"time"

	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"

	"github.com/dmibod/kanban/shared/persistence"

	"github.com/dmibod/kanban/shared/tools/db/mongo"

	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/notify"
	"github.com/dmibod/kanban/process"
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/update"
)

func main() {
	c, cancel := context.WithCancel(context.Background())

	m := mux.ConfigureMux()

	l := createLogger("[MONGO..] ", true)
	f := mongo.CreateFactory(
		"kanban",
		persistence.CreateService(l),
		l)

	boot(&command.Module{Logger: createLogger("[COMMAND] ", true), Mux: m})
	boot(&notify.Module{Logger: createLogger("[NOTIFY.] ", true), Mux: m})

	m.Route("/v1/api/card", func(r chi.Router) {
		router := chi.NewRouter()

		s := services.CreateFactory(l, f)

		monolithic(&query.Module{Logger: createLogger("[QUERY..] ", true), Mux: router, Factory: s})
		monolithic(&update.Module{Logger: createLogger("[UPDATE.] ", true), Mux: router, Factory: s})

		r.Mount("/", router)
	})

	process.Boot(c, createLogger("[PROCESS] ", true))

	mux.StartMux(m, mux.GetPortOrDefault(3000), createLogger("[..MUX..] ", true))

	cancel()

	time.Sleep(time.Second)
}

func boot(b interface{ Boot() }) {
	b.Boot()
}

func monolithic(b interface{ Boot(bool) }) {
	b.Boot(false)
}

func createLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
