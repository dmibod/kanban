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
	BoardRouter chi.Router
	LaneRouter  chi.Router	
	CardRouter  chi.Router
	Factory     *services.Factory
}

// Boot installs handlers to mux
func (m *Module) Boot() {
	if m.Logger == nil {
		m.Logger = console.New(console.WithPrefix("[.QUERY.] "), console.WithDebug(true))
	}

	m.Debugln("starting...")

	CreateBoardAPI(m.Factory.CreateBoardService(), m.Logger).Routes(m.BoardRouter)
	CreateLaneAPI(m.Factory.CreateLaneService(), m.Logger).Routes(m.LaneRouter)
	CreateCardAPI(m.Factory.CreateCardService(), m.Logger).Routes(m.CardRouter)

	m.Debugln("started!")
}
