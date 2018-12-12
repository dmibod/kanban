package notify

import (
	
	"net/http"
	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/msg/nats"
	"github.com/dmibod/kanban/tools/mux"
)

func Boot(m mux.Mux){
	var t msg.Transport = nats.New()

	env := &Env{ Msg: t.Receive("notification") }

	m.Handle("/", http.HandlerFunc(env.ServeHome))
	m.Handle("/ws", http.HandlerFunc(env.ServeWs))
}