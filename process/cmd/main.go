package main

import (
	"encoding/json"
	"log"

	"github.com/dmibod/kanban/kernel"

	"github.com/dmibod/kanban/process"
	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/msg/nats"
)

func main() {
	var t msg.Transport = nats.New()

	in := t.Receive("command")
	out := t.Send("notification")

	for m := range in {
		commands := []process.Command{}
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
			case process.InsertCard: //todo
			case process.UpdateCard: //todo
			case process.RemoveCard: //todo
			case process.ExcludeCard: //todo
			}
			if cnt, ok := ids[id]; ok {
				ids[id] = cnt + 1
			} else {
				ids[id] = 1
			}
		}

		if len(ids) == 0 {
			continue
		}
		
		n, jsonErr := json.Marshal(ids)

		if jsonErr != nil {
			log.Println("Error marshal notifiactions")
		} else {
			out <- n
		}
	}
}
