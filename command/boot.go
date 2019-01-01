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
	l := m.Logger
	if l == nil {
		l = console.New(console.WithPrefix("[COMMAND] "), console.WithDebug(true))
	}

	l.Debugln("starting...")

	CreateAPI(message.CreatePublisher("command"), l).Routes(m.Router)

	l.Debugln("started!")
}
