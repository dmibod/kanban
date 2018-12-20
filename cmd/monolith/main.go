package main

import (
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
	"github.com/dmibod/kanban/shared/tools/mux/http"
	"github.com/dmibod/kanban/update"
)

func main() {
	c, cancel := context.WithCancel(context.Background())

	m := http.New()
	l := createLogger("[MONGO..] ", true)
	s := persistence.CreateService(l)
	f := mongo.CreateFactory(mongo.WithDatabase("kanban"), mongo.WithExecutor(s), mongo.WithLogger(l))

	command.Boot(m, createLogger("[COMMAND] ", true))
	notify.Boot(m, createLogger("[NOTIFY.] ", true))
	query.Boot(m, f, createLogger("[QUERY..] ", true))
	update.Boot(m, f, createLogger("[UPDATE.] ", true))
	process.Boot(c, createLogger("[PROCESS] ", true))

	m.Start()

	cancel()

	time.Sleep(time.Second)
}

func createLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
