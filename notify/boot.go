package notify

import (
	"github.com/dmibod/kanban/shared/tools/log"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/msg"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
	"github.com/dmibod/kanban/shared/tools/mux"
)

func Boot(m mux.Mux, l log.Logger) {

	var t msg.Transport = nats.New()

	env := &Env{ Logger: l, NotificationQueue: t.Receive("notification")}

	m.Get("/", http.HandlerFunc(env.ServeHome))
	m.All("/ws", http.HandlerFunc(env.ServeWs))

	l.Infoln("endpoints registered")
}
