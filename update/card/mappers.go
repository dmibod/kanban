package card

import (
	"github.com/dmibod/kanban/shared/services/card"
)

type cardCreateMapper struct {
}

// PayloadToModel mapping
func (cardCreateMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Card)
	return &card.CreateModel{
		Name:        payload.Name,
		Description: payload.Description,
	}
}

// ModelToPayload mapping
func (cardCreateMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*card.Model)
	return &Card{
		ID:          string(model.ID),
		Name:        model.Name,
		Description: model.Description,
	}
}
