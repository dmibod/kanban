package query

import (
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/log/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
)

// Boot installs Query module http handlers to mux
func Boot(m mux.Mux, f db.RepoFactory) {

	l := logger.New(logger.WithPrefix("[QUERY] "), logger.WithDebug(true))

	m.Get("/get", mux.Handle(&GetCard{Logger: l, Service: CreateCardService(l, CreateCardRepository(f))}))

	l.Infoln("endpoints registered")
}
