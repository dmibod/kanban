package card

import (
	"github.com/dmibod/kanban/shared/domain/event"
)

type aggregate struct {
	Entity
	Repository
	event.Bus
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
func (a *aggregate) Save() error {
	entity := &a.Entity
	if err := a.Repository.Update(entity); err != nil {
		return err
	}

	a.Fire()
	return nil
}
