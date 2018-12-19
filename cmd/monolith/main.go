package main

import (
	"context"
	"time"

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
	f := mongo.CreateDatabaseService(nil).CreateFactory(mongo.WithDatabase("kanban"))

	command.Boot(m)
	notify.Boot(m)
	query.Boot(m, f)
	update.Boot(m, f)
	process.Boot(c)

	m.Start()

	cancel()

	time.Sleep(time.Second)
}
