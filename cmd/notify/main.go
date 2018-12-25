package main

import (
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/notify"
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

	module := notify.Module{Logger: createLogger("[NOTIFY.] ", true), Mux: m, Msg: t}

	module.Boot()

	mux.StartMux(m, mux.GetPortOrDefault(8001), createLogger("[..MUX..] ", true))
}

func createLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
