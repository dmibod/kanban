package update

import (
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/tools/mux"
)

// Boot - adds update module handlers to mux
func Boot(m mux.Mux, f db.Factory, l log.Logger) {

	r := persistence.CreateCardRepository(f)
	s := CreateCardService(l, r)
	h := CreateCreateCardHandler(l, s)

	m.Post("/post", mux.Handle(h))

	l.Infoln("endpoints registered")
}
