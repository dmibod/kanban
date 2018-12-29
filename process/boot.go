package process

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/logger"
)

// Module dependencies
type Module struct {
	context.Context
	logger.Logger
}

func (m *Module) Boot() {
	m.Logger.Debugln("starting...")

	env := CreateHandler(m.Logger)

	go env.Handle(m.Context)

	m.Logger.Debugln("started!")
}
