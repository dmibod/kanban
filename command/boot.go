package command

import (
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/go-chi/chi"
)

// Module dependencies
type Module struct {
	Router chi.Router
	logger.Logger
}

// Boot module
func (m *Module) Boot() {
	if m.Logger == nil {
		m.Logger = console.New(console.WithPrefix("[COMMAND] "), console.WithDebug(true))
	}

	m.Debugln("starting...")

	CreateAPI(message.CreatePublisher("command"), m.Logger).Routes(m.Router)

	m.Debugln("started!")
}
