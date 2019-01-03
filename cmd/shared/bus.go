package shared

import (
	"context"
	"os"

	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/bus"
	"github.com/dmibod/kanban/shared/tools/bus/nats"
	"github.com/dmibod/kanban/shared/tools/logger"
)

const busUrlEnvVar = "BUS_URL"
const busClientEnvVar = "BUS_CLIENT"

// GetNameOrDefault gets name from environment variable or fallbacks to default one
func GetNameOrDefault(defName string) string {
	name := os.Getenv(busClientEnvVar)

	if name == "" {
		return defName
	}

	return name
}

func getBusUrlOrDefault(defUrl string) string {
	url := os.Getenv(busUrlEnvVar)

	if url == "" {
		return defUrl
	}

	return url
}

// StartBus starts bus
func StartBus(ctx context.Context, c string, l logger.Logger) {
	conn := nats.CreateConnection(
		nats.WithName(c),
		nats.WithURL(getBusUrlOrDefault("")),
		nats.WithLogger(l))

	if err := bus.ConnectAndServe(ctx, conn, message.CreateTransport(conn, CreateLogger("[BRK.BUS]", true))); err != nil {
		l.Errorf("bus err: %s\n", err.Error()) // panic if there is an error
		panic(err)
	}
}

// StopBus stops bus
func StopBus() {
	bus.Disconnect()
}
