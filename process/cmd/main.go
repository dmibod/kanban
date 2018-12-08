package main

import (
	"encoding/json"
	"log"
	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/msg/nats"
	"github.com/dmibod/kanban/process"
)

func main(){
	var t msg.Transport = nats.New()

	c := t.Receive("command")

	for m := range c {
		commands := []process.Command{}
		jsonErr := json.Unmarshal(m, &commands)
	
		if jsonErr != nil {
			log.Println("Error parsing json", jsonErr)
			return
		}

		log.Println(commands);
	}
}