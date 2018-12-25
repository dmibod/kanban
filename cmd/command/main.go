package main

import (
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
)

func main() {

	l := console.New(
		console.WithPrefix("[COMMAND] "),
		console.WithDebug(true))

	m := mux.ConfigureMux()

	module := command.Module{Logger: l, Mux: m, Msg: nats.CreateTransport(context.Background(), message.CreateService(l))}

	module.Boot()

	mux.StartMux(m, mux.GetPortOrDefault(8000), l)
}
