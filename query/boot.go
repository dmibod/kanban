package query

import (
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/go-chi/chi"
)

// Module dependencies
type Module struct {
	Mux     *chi.Mux
	Factory *services.Factory
	Logger  logger.Logger
}

// Boot installs handlers to mux
func (m *Module) Boot() {

	api := CreateAPI(m.Logger, m.Factory)
	api.Routes(m.Mux)

	m.Logger.Debugln("endpoints registered")
}
