package command

import (
	"encoding/json"
	"net/http"

	"github.com/dmibod/kanban/shared/message"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
)

// Command declares command type at api level
type Command struct {
	ID      kernel.Id          `json:"id"`
	Type    kernel.CommandType `json:"type"`
	Payload map[string]string  `json:"payload"`
}

// API dependencies
type API struct {
	logger.Logger
	message.Publisher
}

// CreateAPI creates new API instance
func CreateAPI(p message.Publisher, l logger.Logger) *API {
	return &API{
		Logger:    l,
		Publisher: p,
	}
}

// Routes install handlers
func (a *API) Routes(router chi.Router) {
	router.Post("/", a.Post)
}

// Post commands to bus
func (a *API) Post(w http.ResponseWriter, r *http.Request) {
	commands := []Command{}

	err := mux.ParseJSON(r, &commands)
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusBadRequest)
		return
	}

	a.Debugf("commands received: %+v\n", commands)

	m, err := json.Marshal(commands)
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	err = a.Publish(m)
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	res := struct {
		Count   int  `json:"count"`
		Success bool `json:"success"`
	}{len(commands), true}

	render.JSON(w, r, res)
}
