package process

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dmibod/kanban/kernel"
)

type Env struct {
	In  <-chan []byte
	Out chan<- []byte
}

func (e *Env) Handle(c context.Context) {
	for {
		select {
		case <-c.Done():
			log.Println("Existing processor routine")
			return
		case m := <-e.In:
			e.process(m)
		}
	}
}

func (e *Env) process(m []byte) {

	commands := []Command{}
	err := json.Unmarshal(m, &commands)

	if err != nil {
		log.Println("Error parsing json", err)
		return
	}

	log.Println(commands)

	ids := make(map[kernel.Id]int)

	for _, c := range commands {
		id := c.Id
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
		log.Println("Error marshal notifiactions")
	} else {
		e.Out <- n
	}
}