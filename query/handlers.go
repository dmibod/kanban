package query

import (
	"net/http"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/log"
)

// Card maps card to/from json at rest api level
type Card struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetCardService service expected by handler
type GetCardService interface {
	GetCardByID(kernel.Id) (*CardModel, error)
}

// GetCardHandler contains dependencies required by handler
type GetCardHandler struct {
	logger  log.Logger
	service GetCardService
}

// CreateGetCardHandler creates new instance of GetCardHandler
func CreateGetCardHandler(l log.Logger, s GetCardService) *GetCardHandler {
	return &GetCardHandler{
		logger:  l,
		service: s,
	}
}

// Parse parses Api request
func (h *GetCardHandler) Parse(r *http.Request) (interface{}, error) {
	return r.FormValue("id"), nil
}

// Handle handles Api request
func (h *GetCardHandler) Handle(req interface{}) (interface{}, error) {
	model, err := h.service.GetCardByID(kernel.Id(req.(string)))
	if err != nil {
		h.logger.Errorln("error getting card", err)
		return nil, err
	}

	return &Card{
		ID:   string(model.ID),
		Name: model.Name,
	}, nil
}
