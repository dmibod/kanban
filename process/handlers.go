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

type Type int

const (
	UpdateCard Type = iota
	RemoveCard
	ExcludeChild
	InsertBefore
	InsertAfter
	AppendChild
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
	mongo.ContextFactory
	laneService services.LaneService
}

// CreateHandler creates handler
func CreateHandler(p message.Publisher, s message.Subscriber, f mongo.ContextFactory, laneService services.LaneService, l logger.Logger) *Handler {
	return &Handler{
		Logger:         l,
		Publisher:      p,
		Subscriber:     s,
		ContextFactory: f,
		laneService:    laneService,
	}
}

// Handle handles message
func (h *Handler) Handle() {
	h.Subscribe(bus.HandleFunc(h.process))
}

func (h *Handler) process(m []byte) {

	commands := []Command{}

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

	ids := make(map[kernel.Id]int)

	for _, c := range commands {
		id := c.ID
		switch c.Type {
		case InsertBefore: //todo
		case InsertAfter: //todo
		case AppendChild:
			laneID, ok := c.Payload["lane_id"]
			if !ok {
				h.Errorln("lane_id is not found in payload of AppendChild command")
			} else {
				err := h.laneService.AppendChild(ctx, kernel.Id(laneID), kernel.Id(c.ID))
				if err != nil {
					h.Errorln(err)
				}
			}
		case UpdateCard: //todo
		case RemoveCard: //todo
		case ExcludeChild:
			laneID, ok := c.Payload["lane_id"]
			if !ok {
				h.Errorln("lane_id is not found in payload of AppendChild command")
			} else {
				err := h.laneService.ExcludeChild(ctx, kernel.Id(laneID), kernel.Id(c.ID))
				if err != nil {
					h.Errorln(err)
				}
			}
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
		h.Errorln(err)
		return
	}

	err = h.Publish(n)
	if err != nil {
		h.Errorln(err)
		return
	}
}
