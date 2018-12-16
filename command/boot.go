package command

import (
	"log"
	"github.com/dmibod/kanban/tools/msg/nats"
	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/mux"
)

// Boot installs command module handlers to mux
func Boot(m mux.Mux){

	var t msg.Transport = nats.New()

	m.Post("/commands", mux.Handle(&PostCommands{ CommandQueue: t.Send("command") }))

	log.Println("Command module endpoints registered")
}