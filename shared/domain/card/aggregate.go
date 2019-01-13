package card

import (
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
)

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

// Delete card
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

// Create card
func Create(id kernel.ID, registry event.Registry) (*Entity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	if registry == nil {
		return nil, err.ErrInvalidArgument
	}

	entity := Entity{ID: id}

	registry.Register(CreatedEvent{entity})

	return &entity, nil
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
