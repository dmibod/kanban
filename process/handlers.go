package process

import (
	"context"
	"encoding/json"

	"github.com/dmibod/kanban/shared/services"

	"github.com/dmibod/kanban/shared/tools/bus"

	"github.com/dmibod/kanban/shared/message"

	"github.com/dmibod/kanban/shared/kernel"
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
	laneService services.LaneService
}

// CreateHandler creates handler
func CreateHandler(p message.Publisher, s message.Subscriber, laneService services.LaneService, l logger.Logger) *Handler {
	return &Handler{
		Logger:      l,
		Publisher:   p,
		Subscriber:  s,
		laneService: laneService,
	}
}

// Handle handles message
func (h *Handler) Handle(c context.Context) {
	h.Subscribe(bus.HandleFunc(func(m []byte) {
		h.process(c, m)
	}))
}

func (h *Handler) process(ctx context.Context, m []byte) {

	commands := []Command{}

	err := json.Unmarshal(m, &commands)
	if err != nil {
		h.Errorln(err)
		return
	}

	h.Debugln(commands)

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
