package lane

import (
	"github.com/dmibod/kanban/shared/services/lane"
)

// ModelMapper type
type ModelMapper struct {
}

// ModelToPayload mapping
func (ModelMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*lane.Model)
	return &Lane{
		ID:     string(model.ID),
		Name:   model.Name,
		Type:   model.Type,
		Layout: model.Layout,
	}
}

// ModelsToPayload mapping
func (m ModelMapper) ModelsToPayload(models []*lane.Model) []interface{} {
	items := []interface{}{}
	for _, model := range models {
		items = append(items, m.ModelToPayload(model))
	}
	return items
}

// ListModelMapper type
type ListModelMapper struct {
}

// ModelToPayload mapping
func (ListModelMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*lane.ListModel)
	return &Lane{
		ID:     string(model.ID),
		Name:   model.Name,
		Type:   model.Type,
		Layout: model.Layout,
	}
}

// ModelsToPayload mapping
func (m ListModelMapper) ModelsToPayload(models []*lane.ListModel) []interface{} {
	items := []interface{}{}
	for _, model := range models {
		items = append(items, m.ModelToPayload(model))
	}
	return items
}
