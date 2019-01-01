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
	CardRouter chi.Router
	Factory    *services.Factory
}

// Boot installs handlers to mux
func (m *Module) Boot() {
	l := m.Logger
	if l == nil {
		l = console.New(console.WithPrefix("[.QUERY.] "), console.WithDebug(true))
	}

	l.Debugln("starting...")

	CreateCardAPI(m.Factory.CreateCardService(), l).Routes(m.CardRouter)

	l.Debugln("started!")
}
