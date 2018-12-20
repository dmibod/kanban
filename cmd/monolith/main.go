package main

import (
	"github.com/go-chi/chi"
	"context"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"

	"github.com/dmibod/kanban/shared/persistence"

	"github.com/dmibod/kanban/shared/tools/db/mongo"

	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/notify"
	"github.com/dmibod/kanban/process"
	"github.com/dmibod/kanban/query"
	utils "github.com/dmibod/kanban/shared/tools/mux"
	"github.com/dmibod/kanban/update"
)

func main() {
	c, cancel := context.WithCancel(context.Background())

	m := utils.ConfigureMux()

	l := createLogger("[MONGO..] ", true)
	s := persistence.CreateService(l)
	f := mongo.CreateFactory(mongo.WithDatabase("kanban"), mongo.WithExecutor(s), mongo.WithLogger(l))

	boot(&command.Env{Logger: createLogger("[COMMAND] ", true), Mux: m })
	boot(&notify.Env {Logger: createLogger("[NOTIFY.] ", true), Mux: m })

	m.Route("/v1/api/card", func(r chi.Router) {
		router := chi.NewRouter()

		boot(&query.Env  {Logger: createLogger("[QUERY..] ", true), Mux: router, Factory: f })
		boot(&update.Env {Logger: createLogger("[UPDATE.] ", true), Mux: router, Factory: f })
	
		r.Mount("/", router)
	})
	
	process.Boot(c, createLogger("[PROCESS] ", true))
	
	utils.PrintRoutes(createLogger("[MUX....] ", true), m)

	utils.StartMux(m, utils.GetPortOrDefault(3000))

	cancel()

	time.Sleep(time.Second)
}

func boot(b interface{ Boot() }){
	b.Boot()
}

func createLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
