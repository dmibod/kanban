package card

import (
	"github.com/dmibod/kanban/shared/services/card"
)

// CardGetMapper mapper
type CardGetMapper struct {
}

// ModelToPayload mapping
func (CardGetMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*card.Model)
	return &Card{
		ID:   string(model.ID),
		Name: model.Name,
	}
}

// List mapping
func (m CardGetMapper) List(models []*card.Model) []interface{} {
	list := make([]interface{}, len(models))
	for i, model := range models {
		list[i] = model
	}
	return list
}
