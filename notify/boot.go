package notify

import (
	"log"
	
	"net/http"
	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/msg/nats"
	"github.com/dmibod/kanban/tools/mux"
)

func Boot(m mux.Mux){
	var t msg.Transport = nats.New()

	env := &Env{ NotificationQueue: t.Receive("notification") }

	m.Get("/", http.HandlerFunc(env.ServeHome))
	m.Any("/ws", http.HandlerFunc(env.ServeWs))

	log.Println("Notification module endpoints registered")
}