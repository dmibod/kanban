package shared

import (
	"os"
	"os/signal"
)

// GetInterruptChan gets interrupt channel
func GetInterruptChan() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	return ch
}
