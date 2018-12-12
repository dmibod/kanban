package main

import (
	"github.com/dmibod/kanban/tools/mux/http"
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/notify"
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/update"
)

func main() {
	mux := http.New(http.WithPort(3000))

	command.Boot(mux)
	notify.Boot(mux)
	query.Boot(mux)
	update.Boot(mux)

  mux.Start()
}
