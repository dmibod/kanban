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

	children := make([]string, len(model.Children))

	for i, id := range model.Children {
		children[i] = id.String()
	}

	return &Model{
		ID:          string(model.ID),
		Name:        model.Name,
		Type:        model.Type,
		Layout:      model.Layout,
		Description: model.Description,
		Children:    children,
	}
}

// ListModelMapper type
type ListModelMapper struct {
}

// ModelToPayload mapping
func (ListModelMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*lane.ListModel)

	return &ListModel{
		ID:          string(model.ID),
		Name:        model.Name,
		Type:        model.Type,
		Layout:      model.Layout,
		Description: model.Description,
	}
}

// List mapping
func (m ListModelMapper) List(models []*lane.ListModel) []interface{} {
	list := make([]interface{}, len(models))

	for i, model := range models {
		list[i] = model
	}

	return list
}
