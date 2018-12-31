package process

import (
	"context"
	"encoding/json"

	"github.com/dmibod/kanban/shared/tools/bus"

	"github.com/dmibod/kanban/shared/message"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type Type int

const (
	UpdateCard Type = iota
	RemoveCard
	ExcludeCard
	InsertCard
)

type Command struct {
	ID      kernel.Id         `json:"id"`
	Type    Type              `json:"type"`
	Payload map[string]string `json:"payload"`
}

// Handler definition
type Handler struct {
	logger.Logger
	message.Publisher
	message.Subscriber
}

// CreateHandler creates handler
func CreateHandler(p message.Publisher, s message.Subscriber, l logger.Logger) *Handler {
	return &Handler{
		Logger:     l,
		Publisher:  p,
		Subscriber: s,
	}
}

// Handle handles message
func (h *Handler) Handle(c context.Context) {
	h.Subscribe(bus.HandleFunc(h.process))
}

func (h *Handler) process(m []byte) {

	commands := []Command{}

	err := json.Unmarshal(m, &commands)
	if err != nil {
		h.Errorln("error parsing json", err)
		return
	}

	h.Debugln(commands)

	ids := make(map[kernel.Id]int)

	for _, c := range commands {
		id := c.ID
		switch c.Type {
		case InsertCard: //todo
		case UpdateCard: //todo
		case RemoveCard: //todo
		case ExcludeCard: //todo
		}
		if cnt, ok := ids[id]; ok {
			ids[id] = cnt + 1
		} else {
			ids[id] = 1
		}
	}

	if len(ids) == 0 {
		return
	}

	n, err := json.Marshal(ids)
	if err != nil {
		h.Errorln("error marshal notifiactions", err)
		return
	}

	err = h.Publish(n)
	if err != nil {
		h.Errorln("error send notifiactions", err)
		return
	}
}
