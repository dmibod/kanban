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
	if standalone {
		m.standalone()
	} else {
		m.monolithic()
	}
}

func (m *Module) standalone() {

	api := CreateAPI(m.Logger, m.Factory)

	m.Mux.Route("/v1/api/card", func(r chi.Router) {
		router := chi.NewRouter()

		api.Routes(router)

		r.Mount("/", router)
	})

	m.Logger.Debugln("endpoints registered")
}

func (m *Module) monolithic() {

	api := CreateAPI(m.Logger, m.Factory)
	api.Routes(m.Mux)

	m.Logger.Debugln("endpoints registered")
}
