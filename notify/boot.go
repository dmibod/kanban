package notify

import (
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/go-chi/chi"
)

// Module dependencies
type Module struct {
	Mux *chi.Mux
	logger.Logger
}

// Boot installs notify module handlers to mux
func (m *Module) Boot() {
	if m.Logger == nil {
		m.Logger = console.New(console.WithPrefix("[.NOTIF.] "), console.WithDebug(true))
	}

	m.Debugln("starting...")

	api := CreateAPI(m.Logger)

	m.Mux.Route("/v1/api/notify", func(r chi.Router) {
		r.Mount("/", api.Routes())
	})

	m.Debugln("started!")
}
