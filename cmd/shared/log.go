package shared

import (
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// CreateLogger creates new logger
func CreateLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
