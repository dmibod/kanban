package update

import (
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
)

// Boot - adds update module handlers to mux
func Boot(m mux.Mux, f db.Factory, l logger.Logger) {

	r := persistence.CreateCardRepository(f)
	s := services.CreateCardService(l, r)
	h := CreateCreateCardHandler(l, s)

	m.Post("/post", mux.Handle(h))

	l.Infoln("endpoints registered")
}
