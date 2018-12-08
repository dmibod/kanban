package main

import (
	"net/http"
	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/msg/nats"
	"github.com/dmibod/kanban/notify"
)

func main() {
	var t msg.Transport = nats.New()

	env := &notify.Env{ Msg: t.Receive("notification") }

	http.HandleFunc("/", env.ServeHome)
	http.HandleFunc("/ws", env.ServeWs)
	http.ListenAndServe(":3001", nil)
}