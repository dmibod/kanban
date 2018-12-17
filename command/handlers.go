package command

import (
	"encoding/json"
	"net/http"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/tools/mux"
)

// CommandType declares command type
type CommandType int

const (
	UpdateCard CommandType = CommandType(iota)
	RemoveCard
	ExcludeCard
	InsertCard
)

// Command declares command type at api level
type Command struct {
	ID      kernel.Id         `json:"id"`
	Type    CommandType       `json:"type"`
	Payload map[string]string `json:"payload"`
}

// PostCommandHandler holds dependencies
type PostCommandHandler struct {
	Logger       log.Logger
	CommandQueue chan<- []byte
}

// CreatePostCommandHandler creates new PostCommandHandler instance
func CreatePostCommandHandler(l log.Logger, q chan<- []byte) *PostCommandHandler {
	return &PostCommandHandler{
		Logger:       l,
		CommandQueue: q,
	}
}

// Parse parses request
func (h *PostCommandHandler) Parse(r *http.Request) (interface{}, error) {
	commands := []Command{}

	err := mux.JsonRequest(r, &commands)
	if err != nil {
		h.Logger.Errorln("Error parsing json", err)
	}

	return commands, err
}

// Handle handles request
func (h *PostCommandHandler) Handle(req interface{}) (interface{}, error) {
	commands := req.([]Command)

	h.Logger.Debugf("Commands received: %+v\n", commands)

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

	h.Logger.Debugf("Commands sent: %+v\n", len(commands))

	return &res, nil
}
