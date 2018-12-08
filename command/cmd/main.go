package main

import (
	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/msg/nats"
	"github.com/dmibod/kanban/command"
)

func main() {
	var t msg.Transport = nats.New()

	env := &command.Env{ msg: t.Send("command") }

	http.HandleFunc("/commands", env.PostCommands)
	http.ListenAndServe(":3000", nil)
}