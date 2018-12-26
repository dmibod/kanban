package main

import (
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/shared/message"
)

func main() {

	s := message.CreateService("COMMAND", shared.CreateLogger("[BRK.NAT] ", true))

	l := shared.CreateLogger("[COMMAND] ", true)

	t := message.CreateTransport(context.Background(), s, l)

	m := shared.ConfigureMux()

	module := command.Module{Mux: m, Msg: t, Logger: l}

	module.Boot()

	shared.StartMux(m, shared.GetPortOrDefault(8000), shared.CreateLogger("[..MUX..] ", true))
}
