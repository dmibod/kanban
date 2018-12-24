package command

import (
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/msg"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
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

	var t msg.Transport = nats.New()

	api := CreateAPI(m.Logger, t.Send("command"))

	m.Mux.Route("/v1/api/commands", func(r chi.Router) {
		r.Mount("/", api.Routes())
	})

	m.Logger.Debugln("started!")
}
