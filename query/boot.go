package query

import (
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/go-chi/chi"
)

// Module dependencies
type Module struct {
	Card    chi.Router
	Factory *services.Factory
	logger.Logger
}

// Boot installs handlers to mux
func (m *Module) Boot() {
	m.Debugln("starting...")

	CreateCardAPI(m.Logger, m.Factory.CreateCardService()).Routes(m.Card)

	m.Debugln("started!")
}
