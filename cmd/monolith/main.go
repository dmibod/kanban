package main

import (
	"github.com/dmibod/kanban/shared/tools/log/logger"
	"context"
	"time"

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
	s := persistence.CreateService(nil)
	f := mongo.CreateFactory(mongo.WithDatabase("kanban"), mongo.WithExecutor(s))

	command.Boot(m,   logger.New(logger.WithPrefix("[COMMAND] "), logger.WithDebug(true)))
	notify.Boot(m,    logger.New(logger.WithPrefix("[NOTIFY ] "), logger.WithDebug(true)))
	query.Boot(m, f,  logger.New(logger.WithPrefix("[QUERY  ] "), logger.WithDebug(true)))
	update.Boot(m, f, logger.New(logger.WithPrefix("[UPDATE ] "), logger.WithDebug(true)))
	process.Boot(c,   logger.New(logger.WithPrefix("[PROCESS] "), logger.WithDebug(true)))

	m.Start()

	cancel()

	time.Sleep(time.Second)
}
