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

// ModelsToPayload mapping
func (m CardGetMapper) ModelsToPayload(models []*card.Model) []interface{} {
	items := []interface{}{}
	for _, model := range models {
		items = append(items, m.ModelToPayload(model))
	}
	return items
}
