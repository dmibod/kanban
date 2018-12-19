package command

import (
	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
	"github.com/dmibod/kanban/shared/tools/msg"
	"github.com/dmibod/kanban/shared/tools/mux"
)

// Boot installs command module handlers to mux
func Boot(m mux.Mux, l log.Logger){

	var t msg.Transport = nats.New()

	m.Post("/commands", mux.Handle(CreatePostCommandHandler(l, t.Send("command"))))

	l.Debugln("endpoints registered")
}