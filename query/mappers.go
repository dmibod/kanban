package query

import (
	"github.com/dmibod/kanban/shared/services"
)

// BoardGetMapper mapper
type BoardGetMapper struct {
}

// ModelToPayload mapping
func (BoardGetMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.BoardModel)
	return &Board{
		ID:     string(model.ID),
		Name:   model.Name,
		Layout: model.Layout,
		Owner:  model.Owner,
		Shared: model.Shared,
	}
}

// ModelsToPayload mapping
func (m BoardGetMapper) ModelsToPayload(models []*services.BoardModel) []interface{} {
	items := []interface{}{}
	for _, model := range models {
		items = append(items, m.ModelToPayload(model))
	}
	return items
}

// LaneGetMapper mapper
type LaneGetMapper struct {
}

// ModelToPayload mapping
func (LaneGetMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.LaneModel)
	return &Lane{
		ID:     string(model.ID),
		Name:   model.Name,
		Type:   model.Type,
		Layout: model.Layout,
	}
}

// ModelsToPayload mapping
func (m LaneGetMapper) ModelsToPayload(models []*services.LaneModel) []interface{} {
	items := []interface{}{}
	for _, model := range models {
		items = append(items, m.ModelToPayload(model))
	}
	return items
}

// CardGetMapper mapper
type CardGetMapper struct {
}

// ModelToPayload mapping
func (CardGetMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.CardModel)
	return &Card{
		ID:   string(model.ID),
		Name: model.Name,
	}
}

// ModelsToPayload mapping
func (m CardGetMapper) ModelsToPayload(models []*services.CardModel) []interface{} {
	items := []interface{}{}
	for _, model := range models {
		items = append(items, m.ModelToPayload(model))
	}
	return items
}
