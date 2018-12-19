package query

import (
	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/dmibod/kanban/shared/persistence"
)

// Boot installs Query module http handlers to mux
func Boot(m mux.Mux, f db.Factory, l log.Logger) {

	r := persistence.CreateCardRepository(f)
	s := CreateCardService(l, r)
	h := CreateGetCardHandler(l, s)

	m.Get("/get", mux.Handle(h))

	l.Debugln("endpoints registered")
}
