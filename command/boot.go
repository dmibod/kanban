package command

import (
	"github.com/dmibod/kanban/shared/tools/msg/nats"
	"github.com/dmibod/kanban/shared/tools/msg"
	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/dmibod/kanban/shared/tools/log/logger"
)

// Boot installs command module handlers to mux
func Boot(m mux.Mux){

	l := logger.New(logger.WithPrefix("[COMMAND] "), logger.WithDebug(true))

	var t msg.Transport = nats.New()

	m.Post("/commands", mux.Handle(CreatePostCommandHandler(l, t.Send("command"))))

	l.Debugln("endpoints registered")
}