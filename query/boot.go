package query

import (
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/dmibod/kanban/shared/persistence"
)

// Boot installs Query module http handlers to mux
func Boot(m mux.Mux, f db.Factory, l logger.Logger) {

	r := persistence.CreateCardRepository(f)
	s := services.CreateCardService(l, r)
	h := CreateGetCardHandler(l, s)

	m.Get("/get", mux.Handle(h))

	l.Debugln("endpoints registered")
}
