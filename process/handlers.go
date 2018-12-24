package process

import (
	"context"
	"encoding/json"

	"github.com/dmibod/kanban/shared/tools/msg"

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

type Env struct {
	Logger   logger.Logger
	Sender   msg.Sender
	Receiver msg.Receiver
	Queue    chan []byte
}

func CreateHandler(l logger.Logger, s msg.Sender, r msg.Receiver) *Env {
	q := make(chan []byte)
	err := r.Receive("", func(msg []byte) {
		q <- msg
	})
	if err != nil {
		l.Errorln("error subscribe queue", err)
	}

	return &Env{Logger: l, Sender: s, Receiver: r, Queue: q}
}

func (e *Env) Handle(c context.Context) {
	for {
		select {
		case <-c.Done():
			e.Logger.Debugln("Exiting processor routine")
			return
		case m := <-e.Queue:
			e.process(m)
		}
	}
}

func (e *Env) process(m []byte) {

	commands := []Command{}
	err := json.Unmarshal(m, &commands)

	if err != nil {
		e.Logger.Errorln("error parsing json", err)
		return
	}

	e.Logger.Debugln(commands)

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

	n, jsonErr := json.Marshal(ids)

	if jsonErr != nil {
		e.Logger.Errorln("error marshal notifiactions")
		return
	}

	sendErr := e.Sender.Send(n)
	if sendErr != nil {
		e.Logger.Errorln("error send notifiactions", sendErr)
		return
	}
}
