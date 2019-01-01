package process

import (
	"context"

	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
)

// Module dependencies
type Module struct {
	context.Context
	logger.Logger
	Factory *services.Factory
}

// Boot installs module handlers to bus
func (m *Module) Boot() {
	if m.Logger == nil {
		m.Logger = console.New(console.WithPrefix("[PROCESS] "), console.WithDebug(true))
	}

	m.Debugln("starting...")

	h := CreateHandler(
		message.CreatePublisher("notification"),
		message.CreateSubscriber("command"),
		m.Factory.CreateLaneService(),
		m.Logger)

	h.Handle(m.Context)

	m.Debugln("started!")
}
