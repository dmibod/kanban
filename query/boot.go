package query

import (
	"github.com/dmibod/kanban/tools/db"
	"github.com/dmibod/kanban/tools/log/logger"
	"github.com/dmibod/kanban/tools/mux"
)

// Boot installs Query module http handlers to mux
func Boot(m mux.Mux, f db.RepoFactory) {

	l := logger.New(logger.WithPrefix("Query"))

	m.Get("/get", mux.Handle(&GetCard{Logger: l, Service: CreateCardService(l, CreateCardRepository(f))}))

	l.Infoln("Query module endpoints registered")
}
