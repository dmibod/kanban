package update

import (
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/go-chi/chi"
)

// Module dependencies
type Module struct {
	Card    chi.Router
	Board   chi.Router
	Factory *services.Factory
	Logger  logger.Logger
}

// Boot installs handlers to mux
func (m *Module) Boot() {
	m.Logger.Debugln("starting...")

	CreateCardAPI(m.Logger, m.Factory).Routes(m.Card)
	CreateBoardAPI(m.Logger, m.Factory).Routes(m.Board)

	m.Logger.Debugln("started!")
}
