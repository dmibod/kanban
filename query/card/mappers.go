package card

import (
	"github.com/dmibod/kanban/shared/services/card"
)

// ModelMapper mapper
type ModelMapper struct {
}

// ModelToPayload mapping
func (ModelMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*card.Model)
	return &Model{
		ID:          string(model.ID),
		Name:        model.Name,
		Description: model.Description,
	}
}

// List mapping
func (m ModelMapper) List(models []*card.Model) []interface{} {
	list := make([]interface{}, len(models))
	for i, model := range models {
		list[i] = model
	}
	return list
}
