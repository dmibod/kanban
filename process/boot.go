package process

import (
	"context"
	"github.com/dmibod/kanban/shared/message"

	"github.com/dmibod/kanban/shared/tools/logger"
)

// Module dependencies
type Module struct {
	context.Context
	logger.Logger
}

// Boot installs module handlers to bus
func (m *Module) Boot() {
	m.Debugln("starting...")

	h := CreateHandler(
		message.CreatePublisher("notification"),
		message.CreateSubscriber("command"),
		m.Logger)

	go h.Handle(m.Context)

	m.Debugln("started!")
}
