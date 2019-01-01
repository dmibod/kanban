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
	CardRouter  chi.Router
	BoardRouter chi.Router
	Factory     *services.Factory
}

// Boot installs handlers to mux
func (m *Module) Boot() {
	if m.Logger == nil {
		m.Logger = console.New(console.WithPrefix("[.QUERY.] "), console.WithDebug(true))
	}

	m.Debugln("starting...")

	CreateCardAPI(m.Factory.CreateCardService(), m.Logger).Routes(m.CardRouter)
	CreateBoardAPI(m.Factory.CreateBoardService(), m.Logger).Routes(m.BoardRouter)

	m.Debugln("started!")
}
