package main

import (
	"github.com/dmibod/kanban/shared/tools/log/logger"
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/mux/http"
)

func main() {
	l := logger.New(logger.WithPrefix("[QUERY..] "), logger.WithDebug(true))
	m := http.New(http.WithPort(http.GetPortOrDefault(3002)))
	s := persistence.CreateService(l)
	f := mongo.CreateFactory(mongo.WithDatabase("kanban"), mongo.WithExecutor(s), mongo.WithLogger(l))

	query.Boot(m, f, l)

	m.Start()
}
