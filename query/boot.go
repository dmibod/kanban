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
func (m *Module) Boot(standalone bool) {
	m.Logger.Debugln("starting...")

	if standalone {
		m.standalone()
	} else {
		m.monolithic()
	}

	m.Logger.Debugln("started!")
}

func (m *Module) standalone() {

	api := CreateAPI(m.Logger, m.Factory)

	m.Mux.Route("/v1/api/card", func(r chi.Router) {
		router := chi.NewRouter()

		api.Routes(router)

		r.Mount("/", router)
	})
}

func (m *Module) monolithic() {

	CreateAPI(m.Logger, m.Factory).Routes(m.Mux)
}
