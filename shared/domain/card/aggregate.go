package card

import (
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
)

// Create card
func Create(id kernel.ID, bus event.Bus) (*Entity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	if bus == nil {
		return nil, err.ErrInvalidArgument
	}

	entity := Entity{ID: id}

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

// Delete card
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

// Save changes
func (a *aggregate) Save() {
	a.Fire()
}
