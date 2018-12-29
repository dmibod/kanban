package command

import (
	"encoding/json"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/bus"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
)

// Type declares command type
type Type int

const (
	UpdateCard Type = Type(iota)
	RemoveCard
	ExcludeCard
	InsertCard
)

// Command declares command type at api level
type Command struct {
	ID      kernel.Id         `json:"id"`
	Type    Type              `json:"type"`
	Payload map[string]string `json:"payload"`
}

// API holds dependencies required by handlers
type API struct {
	logger logger.Logger
}

// CreateAPI creates new API instance
func CreateAPI(l logger.Logger) *API {
	return &API{
		logger: l,
	}
}

// Routes export API router
func (a *API) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", a.Post)
	return router
}

// Post - posts commands to queue
func (a *API) Post(w http.ResponseWriter, r *http.Request) {
	commands := []Command{}

	err := mux.JsonRequest(r, &commands)
	if err != nil {
		a.logger.Errorln("error parsing json", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	a.logger.Debugf("commands received: %+v\n", commands)

	m, err := json.Marshal(commands)
	if err != nil {
		a.logger.Errorln("error marshalling commands", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	err = bus.Publish("command", m)
	if err != nil {
		a.logger.Errorln("error sending commands", err)
		mux.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	a.logger.Debugf("commands sent: %+v\n", len(commands))

	res := struct {
		Count   int  `json:"count"`
		Success bool `json:"success"`
	}{len(commands), true}

	render.JSON(w, r, res)
}
