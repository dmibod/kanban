package process

import (
	"context"
	"encoding/json"

	"github.com/dmibod/kanban/shared/services"

	"github.com/dmibod/kanban/shared/tools/bus"

	"github.com/dmibod/kanban/shared/message"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// Handler definition
type Handler struct {
	logger.Logger
	message.Subscriber
	mongo.ContextFactory
	services.CommandService
}

// CreateHandler creates handler
func CreateHandler(s message.Subscriber, f mongo.ContextFactory, service services.CommandService, l logger.Logger) *Handler {
	return &Handler{
		Logger:         l,
		Subscriber:     s,
		ContextFactory: f,
		CommandService: service,
	}
}

// Handle handles message
func (h *Handler) Handle() {
	h.Subscribe(bus.HandleFunc(h.process))
}

func (h *Handler) process(m []byte) {

	commands := []kernel.Command{}

	err := json.Unmarshal(m, &commands)
	if err != nil {
		h.Errorln(err)
		return
	}

	h.Debugln(commands)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx, err = h.ContextFactory.Context(ctx)
	if err != nil {
		h.Errorln(err)
		return
	}

	for _, c := range commands {
		err = h.CommandService.Execute(ctx, c)
		if err != nil {
			h.Errorln(err)
		}
	}
}
