package process

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/logger"

	"github.com/dmibod/kanban/shared/tools/msg"
)

// Module dependencies
type Module struct {
	Ctx    context.Context
	Msg    msg.Transport
	Logger logger.Logger
}

func (m *Module) Boot() {
	m.Logger.Debugln("starting...")

	env := CreateHandler(m.Logger, m.Msg.CreateSender("notification"), m.Msg.CreateReceiver("command"))

	go env.Handle(m.Ctx)

	m.Logger.Debugln("started!")
}
