package shared

import (
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
)

const debug = true

// CreateLogger creates new logger
func CreateLogger(prefix string) logger.Logger {
	return console.New(
		console.WithPrefix(prefix),
		console.WithDebug(debug))
}
