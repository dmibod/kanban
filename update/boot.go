package update

import (
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/update/board"
	"github.com/dmibod/kanban/update/card"
	"github.com/dmibod/kanban/update/lane"
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
		m.Logger = console.New(console.WithPrefix("[.UPDAT.] "), console.WithDebug(true))
	}

	m.Debugln("starting...")

	board.CreateAPI(m.ServiceFactory.CreateBoardService(), m.Logger).Routes(m.BoardRouter)
	lane.CreateAPI(m.ServiceFactory.CreateLaneService(), m.Logger).Routes(m.LaneRouter)
	card.CreateAPI(m.ServiceFactory.CreateCardService(), m.Logger).Routes(m.CardRouter)

	m.Debugln("started!")
}
