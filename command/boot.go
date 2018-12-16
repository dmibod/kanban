package command

import (
	"github.com/dmibod/kanban/tools/msg/nats"
	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/mux"
	"github.com/dmibod/kanban/tools/log/logger"
)

// Boot installs command module handlers to mux
func Boot(m mux.Mux){

	l := logger.New(logger.WithPrefix("[COMMAND] "), logger.WithDebug(true))

	var t msg.Transport = nats.New()

	m.Post("/commands", mux.Handle(&PostCommands{ CommandQueue: t.Send("command") }))

	l.Infoln("endpoints registered")
}