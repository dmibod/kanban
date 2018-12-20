package update

import (
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/kernel"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/mux"
)

// Card maps card to/from json at rest api level
type Card struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// CardService service expected by handler
type CardService interface{
	CreateCard(*services.CardPayload) (kernel.Id, error)
}

// CreateCardHandler contains dependencies required by handler
type CreateCardHandler struct {
	logger  logger.Logger
	service CardService
}

// CreateCreateCardHandler creates new CreateCardHandler instance
func CreateCreateCardHandler(l logger.Logger, s CardService) *CreateCardHandler {
	return &CreateCardHandler{
		logger:  l,
		service: s,
	}
}

// Parse parse request
func (h *CreateCardHandler) Parse(r *http.Request) (interface{}, error) {
	card := &Card{}

	err := mux.JsonRequest(r, card)
	if err != nil {
		h.logger.Errorln("error parsing json", err)
	}

	return card, err
}

// Handle handles request
func (h *CreateCardHandler) Handle(req interface{}) (interface{}, error) {
	card := req.(*Card)

	id, err := h.service.CreateCard(&services.CardPayload{Name: card.Name})
	if err != nil {
		h.logger.Errorln("error inserting document", err)
		return nil, err
	}

	res := struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}{string(id), true}

	return &res, nil
}
