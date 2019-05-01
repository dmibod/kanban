package query

import (
	"github.com/dmibod/kanban/query/board"
	"github.com/dmibod/kanban/query/card"
	"github.com/dmibod/kanban/query/lane"
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

	board.CreateAPI(b, m.Logger).Routes(m.BoardRouter)
	lane.CreateAPI(l, m.Logger).Routes(m.LaneRouter)
	card.CreateAPI(c, m.Logger).Routes(m.CardRouter)

	m.Debugln("started!")
}
