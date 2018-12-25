package main

import (
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
)

func main() {

	m := mux.ConfigureMux()

	t := nats.CreateTransport(
		context.Background(),
		message.CreateService(createLogger("[BRK.NAT] ", true)))

	module := command.Module{Logger: createLogger("[COMMAND] ", true), Mux: m, Msg: t}

	module.Boot()

	mux.StartMux(m, mux.GetPortOrDefault(8000), createLogger("[..MUX..] ", true))
}

func createLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
