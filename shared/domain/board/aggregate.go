package board

import (
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
)

// Create board
func Create(id kernel.ID, owner string, bus event.Bus) (*Entity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	if owner == "" || bus == nil {
		return nil, err.ErrInvalidArgument
	}

	entity := Entity{
		ID:       id,
		Owner:    owner,
		Layout:   kernel.VLayout,
		Shared:   false,
		Children: []kernel.ID{},
	}

	bus.Register(CreatedEvent{entity})
	bus.Fire()

	return &entity, nil
}

// New aggregate
func New(entity Entity, bus event.Bus) (Aggregate, error) {
	if !entity.ID.IsValid() {
		return nil, err.ErrInvalidID
	}

	if bus == nil {
		return nil, err.ErrInvalidArgument
	}

	return &aggregate{
		Entity: entity,
		Bus:    bus,
	}, nil
}

// Delete board
func Delete(entity Entity, bus event.Bus) error {
	if !entity.ID.IsValid() {
		return err.ErrInvalidID
	}

	if bus == nil {
		return err.ErrInvalidArgument
	}

	bus.Register(DeletedEvent{entity})
	bus.Fire()

	return nil
}

type aggregate struct {
	event.Bus
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
	if a.Entity.Layout == value {
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

// Shared update
func (a *aggregate) Shared(value bool) error {
	if a.Entity.Shared == value {
		return nil
	}

	event := SharedChangedEvent{
		ID:       a.Entity.ID,
		OldValue: a.Entity.Shared,
		NewValue: value,
	}

	a.Entity.Shared = value

	a.Register(event)

	return nil
}

// AppendChild to board
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

// RemoveChild to board
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

// Save changes
func (a *aggregate) Save() {
	a.Fire()
}

func (a *aggregate) findChildIndex(id kernel.ID) int {
	for i, childID := range a.Entity.Children {
		if childID == id {
			return i
		}
	}
	return -1
}
