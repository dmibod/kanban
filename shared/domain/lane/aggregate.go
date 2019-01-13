package lane

import (
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
)

// Create lane
func Create(id kernel.ID, kind string, registry event.Registry) (*Entity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	if registry == nil {
		return nil, err.ErrInvalidArgument
	}

	if kind != kernel.LKind && kind != kernel.CKind {
		return nil, err.ErrInvalidArgument
	}

	layout := kernel.VLayout

	if kind == kernel.CKind {
		layout = ""
	}

	entity := Entity{
		ID:       id,
		Kind:     kind,
		Layout:   layout,
		Children: []kernel.ID{},
	}

	registry.Register(CreatedEvent{entity})

	return &entity, nil
}

// New aggregate
func New(entity Entity, registry event.Registry) (Aggregate, error) {
	if !entity.ID.IsValid() {
		return nil, err.ErrInvalidID
	}

	if registry == nil {
		return nil, err.ErrInvalidArgument
	}

	return &aggregate{
		Entity:   entity,
		Registry: registry,
	}, nil
}

// Delete lane
func Delete(entity Entity, registry event.Registry) error {
	if !entity.ID.IsValid() {
		return err.ErrInvalidID
	}

	if registry == nil {
		return err.ErrInvalidArgument
	}

	registry.Register(DeletedEvent{entity})

	return nil
}

type aggregate struct {
	event.Registry
	Entity
}

// Root entity
func (a *aggregate) Root() Entity {
	return a.Entity
}

// Name update
func (a *aggregate) Name(value string) error {
	if a.Entity.Name == value {
		return nil
	}

	event := NameChangedEvent{
		ID:       a.Entity.ID,
		OldValue: a.Entity.Name,
		NewValue: value,
	}

	a.Entity.Name = value

	a.Register(event)

	return nil
}

// Description update
func (a *aggregate) Description(value string) error {
	if a.Entity.Description == value {
		return nil
	}

	event := DescriptionChangedEvent{
		ID:       a.Entity.ID,
		OldValue: a.Entity.Description,
		NewValue: value,
	}

	a.Entity.Description = value

	a.Register(event)

	return nil
}

// Layout update
func (a *aggregate) Layout(value string) error {
	if a.Entity.Layout == value || a.Entity.Kind == kernel.CKind {
		return nil
	}

	if value == kernel.VLayout || value == kernel.HLayout {
		event := LayoutChangedEvent{
			ID:       a.Entity.ID,
			OldValue: a.Entity.Layout,
			NewValue: value,
		}

		a.Entity.Layout = value

		a.Register(event)

		return nil
	}

	return err.ErrInvalidArgument
}

// AppendChild to lane
func (a *aggregate) AppendChild(id kernel.ID) error {
	if !id.IsValid() {
		return err.ErrInvalidID
	}

	i := a.findChildIndex(id)
	if i < 0 {
		event := ChildAppendedEvent{
			ID:      a.Entity.ID,
			ChildID: id,
		}

		a.Entity.Children = append(a.Entity.Children, id)

		a.Register(event)
	}

	return nil
}

// RemoveChild to lane
func (a *aggregate) RemoveChild(id kernel.ID) error {
	if !id.IsValid() {
		return err.ErrInvalidID
	}

	i := a.findChildIndex(id)
	if i >= 0 {
		event := ChildRemovedEvent{
			ID:      a.Entity.ID,
			ChildID: a.Entity.Children[i],
		}

		a.Entity.Children = append(a.Entity.Children[:i], a.Entity.Children[i+1:]...)

		a.Register(event)
	}

	return nil
}

func (a *aggregate) findChildIndex(id kernel.ID) int {
	for i, childID := range a.Entity.Children {
		if childID == id {
			return i
		}
	}
	return -1
}
