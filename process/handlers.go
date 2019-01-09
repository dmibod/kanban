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
	message.Publisher
	message.Subscriber
	mongo.ContextFactory
	services.CommandService
}

// CreateHandler creates handler
func CreateHandler(p message.Publisher, s message.Subscriber, f mongo.ContextFactory, service services.CommandService, l logger.Logger) *Handler {
	return &Handler{
		Logger:         l,
		Publisher:      p,
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

	notifications := []*kernel.Notification{}

	for _, c := range commands {
		err = h.CommandService.Execute(ctx, c)
		if err != nil {
			h.Errorln(err)
		} else if n := makeNotification(c); n != nil {
			notifications = append(notifications, n)
		}
	}

	if len(notifications) == 0 {
		return
	}

	n, err := json.Marshal(notifications)
	if err != nil {
		h.Errorln(err)
		return
	}

	err = h.Publish(n)
	if err != nil {
		h.Errorln(err)
	}
}

func makeNotification(command kernel.Command) *kernel.Notification {
	if command.Type == kernel.LayoutBoard {
		return &kernel.Notification{
			Context: command.ID,
			ID:      command.ID,
			Type:    kernel.RefreshBoard,
		}
	}

	return nil
}
