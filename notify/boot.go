package notify

import (
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/go-chi/chi"
)

// Module dependencies
type Module struct {
	chi.Router
	logger.Logger
}

// Boot installs notify module handlers to mux
func (m *Module) Boot() {
	if m.Logger == nil {
		m.Logger = console.New(console.WithPrefix("[.NOTIF.] "), console.WithDebug(true))
	}

	m.Debugln("starting...")

	CreateAPI(message.CreateSubscriber("notification"), m.Logger).Routes(m.Router)

	m.Debugln("started!")
}
