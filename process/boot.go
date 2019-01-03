package process

import (
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
)

// Module dependencies
type Module struct {
	mongo.ContextFactory
	logger.Logger
	ServiceFactory *services.ServiceFactory
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
		m.ContextFactory,
		m.ServiceFactory.CreateLaneService(),
		m.Logger)

	h.Handle()

	m.Debugln("started!")
}
