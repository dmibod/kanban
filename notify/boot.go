package notify

import (
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/msg"
)

// Module dependencies
type Module struct {
	Mux *chi.Mux
	msg.Transport
	logger.Logger
}

// Boot installs notify module handlers to mux
func (m *Module) Boot() {
	m.Debugln("starting...")

	api := CreateAPI(m.Logger, m.Subscriber("notification"))

	m.Mux.Route("/v1/api/notify", func(r chi.Router) {
		r.Mount("/", api.Routes())
	})

	m.Debugln("started!")
}
