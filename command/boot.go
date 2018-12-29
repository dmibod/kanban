package command

import (
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/go-chi/chi"
)

// Module dependencies
type Module struct {
	Mux    *chi.Mux
	Logger logger.Logger
}

// Boot installs command module handlers to mux
func (m *Module) Boot() {
	m.Logger.Debugln("starting...")

	api := CreateAPI(m.Logger)

	m.Mux.Route("/v1/api/commands", func(r chi.Router) {
		r.Mount("/", api.Routes())
	})

	m.Logger.Debugln("started!")
}
