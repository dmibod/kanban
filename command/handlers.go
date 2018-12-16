package command

import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/dmibod/kanban/tools/mux"
)

// PostCommands containes dependencies required by handler
type PostCommands struct {
	CommandQueue chan<- []byte
}

// Parse parses request
func (h *PostCommands) Parse(r *http.Request) (interface{}, error){
	commands := []Command{}
	err := mux.JsonRequest(r, &commands)
	if err != nil {
		log.Println("Error parsing json", err)
	}
	return commands, err
}

// Handle handles request
func (h *PostCommands) Handle(req interface{}) (interface{}, error){
	commands := req.([]Command)

	log.Printf("Commands received: %+v\n", commands);

	m, err := json.Marshal(commands)
	if err != nil {
		log.Println("Error marshalling commands", err)
		return nil, err
	}

	h.CommandQueue <- m

	res := struct {
		Count int `json:"count"`
		Success bool `json:"success"`
	}{len(commands),true}

	log.Printf("Commands sent: %+v\n", len(commands))

	return &res, nil
}
