package notify

import (
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/msg"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
)

// Module dependencies
type Module struct {
	Mux    *chi.Mux
	Logger logger.Logger
}

// Boot installs notify module handlers to mux
func (m *Module) Boot() {
	m.Logger.Debugln("starting...")

	var t msg.Transport = nats.New()

	api := CreateAPI(m.Logger, t.Receive("notification"))

	m.Mux.Route("/v1/api/notify", func(r chi.Router) {
		r.Mount("/", api.Routes())
	})

	m.Logger.Debugln("started!")
}
