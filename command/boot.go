package command

import (
	"github.com/dmibod/kanban/tools/msg/nats"
	"github.com/dmibod/kanban/tools/msg"
	"net/http"
	"github.com/dmibod/kanban/tools/mux"
)

func Boot(m mux.Mux){

	var t msg.Transport = nats.New()

	env := &Env{ CommandQueue: t.Send("command") }

	m.Handle("/commands", http.HandlerFunc(env.PostCommands))
}