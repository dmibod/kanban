package query

import (
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/go-chi/chi"
)

// Module dependencies
type Module struct {
	logger.Logger
	BoardRouter    chi.Router
	LaneRouter     chi.Router
	CardRouter     chi.Router
	ServiceFactory *services.ServiceFactory
}

// Boot installs handlers to mux
func (m *Module) Boot() {
	if m.Logger == nil {
		m.Logger = console.New(console.WithPrefix("[.QUERY.] "), console.WithDebug(true))
	}

	m.Debugln("starting...")

	b := m.ServiceFactory.CreateBoardService()
	l := m.ServiceFactory.CreateLaneService()
	c := m.ServiceFactory.CreateCardService()

	CreateBoardAPI(b, l, m.Logger).Routes(m.BoardRouter)
	CreateLaneAPI(l, c, m.Logger).Routes(m.LaneRouter)
	CreateCardAPI(c, m.Logger).Routes(m.CardRouter)

	m.Debugln("started!")
}
