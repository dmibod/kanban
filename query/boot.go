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
func (m *Module) Boot(standalone bool) {
	m.Debugln("starting...")

	CreateCardAPI(m.Logger, m.Factory).Routes(m.Card)

	m.Debugln("started!")
}
