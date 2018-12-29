package shared

import (
	"os"
	"context"

	buz "github.com/dmibod/kanban/shared/bus"
	"github.com/dmibod/kanban/shared/tools/bus"
	"github.com/dmibod/kanban/shared/tools/bus/nats"
	"github.com/dmibod/kanban/shared/tools/logger"
)

const busClientEnvVar = "BUS_CLIENT"

// GetNameOrDefault gets name from environment variable or fallbacks to default one
func GetNameOrDefault(defName string) string {
	name := os.Getenv(busClientEnvVar)

	if name == "" {
		return defName
	}

	return name
}

// StartBus starts bus
func StartBus(ctx context.Context, c string, l logger.Logger) {
	conn := nats.CreateConnection(
		nats.WithName(c),
		nats.WithLogger(l))

	if err := bus.ConnectAndServe(ctx, conn, buz.CreateTransport(conn, l)); err != nil {
		l.Errorf("Mux err: %s\n", err.Error()) // panic if there is an error
		panic(err)
	}
}
