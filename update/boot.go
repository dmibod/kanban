package update

import (
	"github.com/dmibod/kanban/tools/db"
	"github.com/dmibod/kanban/tools/mux"
	"github.com/dmibod/kanban/tools/log/logger"
)

// Boot - adds update module handlers to mux
func Boot(m mux.Mux, f db.RepoFactory){

	l := logger.New(logger.WithPrefix("[UPDATE] "))

	instance := func() interface{} {
		return &Card{}
	}

	m.Post("/post", mux.Handle(&CreateCard{ Repository: f.Create("cards", instance) }))

	l.Infoln("endpoints registered")
}