package command

import (
	"encoding/json"
	"net/http"

	"github.com/dmibod/kanban/tools/log"
	"github.com/dmibod/kanban/tools/mux"
)

// PostCommands containes dependencies required by handler
type PostCommands struct {
	Logger       log.Logger
	CommandQueue chan<- []byte
}

// Parse parses request
func (h *PostCommands) Parse(r *http.Request) (interface{}, error) {
	commands := []Command{}
	err := mux.JsonRequest(r, &commands)
	if err != nil {
		h.Logger.Errorln("Error parsing json", err)
	}
	return commands, err
}

// Handle handles request
func (h *PostCommands) Handle(req interface{}) (interface{}, error) {
	commands := req.([]Command)

	h.Logger.Infof("Commands received: %+v\n", commands)

	m, err := json.Marshal(commands)
	if err != nil {
		h.Logger.Errorln("Error marshalling commands", err)
		return nil, err
	}

	h.CommandQueue <- m

	res := struct {
		Count   int  `json:"count"`
		Success bool `json:"success"`
	}{len(commands), true}

	h.Logger.Infof("Commands sent: %+v\n", len(commands))

	return &res, nil
}
