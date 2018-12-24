package main

import (
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/notify"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
)

func main() {

	l := console.New(
		console.WithPrefix("[NOTIFY.] "),
		console.WithDebug(true))

	m := mux.ConfigureMux()

	module := notify.Module{Logger: l, Mux: m, Msg: nats.CreateTransport(context.Background(), message.CreateService(l))}

	module.Boot()

	mux.StartMux(m, mux.GetPortOrDefault(3001), l)
}
