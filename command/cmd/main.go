package main

import (
	"net/http"
	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/msg/nats"
	"github.com/dmibod/kanban/command"
)

func main() {
	var t msg.Transport = nats.New()

	env := &command.Env{ CommandQueue: t.Send("command") }

	http.HandleFunc("/commands", env.PostCommands)
	
	http.ListenAndServe(":3000", nil)
}