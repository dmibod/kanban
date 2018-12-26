package main

import (
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/notify"
	"github.com/dmibod/kanban/shared/message"
)

func main() {

	s := message.CreateService("NOTIFY", shared.CreateLogger("[BRK.NAT] ", true))

	l := shared.CreateLogger("[NOTIFY.] ", true)

	t := message.CreateTransport(context.Background(), s, l)

	m := shared.ConfigureMux()

	module := notify.Module{Mux: m, Transport: t, Logger: l}

	module.Boot()

	shared.StartMux(m, shared.GetPortOrDefault(8001), shared.CreateLogger("[..MUX..] ", true))
}
