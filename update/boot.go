package update

import (
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/dmibod/kanban/shared/tools/log/logger"
)

// Boot - adds update module handlers to mux
func Boot(m mux.Mux, f db.RepoFactory){

	l := logger.New(logger.WithPrefix("[UPDATE] "), logger.WithDebug(true))

	instance := func() interface{} {
		return &Card{}
	}

	m.Post("/post", mux.Handle(&CreateCard{ Repository: f.Create("cards", instance) }))

	l.Infoln("endpoints registered")
}