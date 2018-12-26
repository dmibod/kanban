package main

import (
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
)

func main() {

	s := message.CreateService("COMMAND", createLogger("[BRK.NAT] ", true))

	l := createLogger("[COMMAND] ", true)

	t := message.CreateTransport(context.Background(), s, l)

	m := mux.ConfigureMux()

	module := command.Module{Mux: m, Msg: t, Logger: l}

	module.Boot()

	mux.StartMux(m, mux.GetPortOrDefault(8000), createLogger("[..MUX..] ", true))
}

func createLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
