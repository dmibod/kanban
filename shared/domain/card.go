package domain

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// CardEntity type
type CardEntity struct {
	ID          kernel.Id
	Name        string
	Description string
}

// CardAggregate interface
type CardAggregate interface {
	Saver
}

type cardAggregate struct {
	Repository
	EventRegistry
	id          kernel.Id
	name        string
	description string
}

func (a *cardAggregate) getEntity() CardEntity {
	return CardEntity{
		ID:          a.id,
		Name:        a.name,
		Description: a.description,
	}
}

// Save changes
func (a *cardAggregate) Save() error {
	id, err := a.Repository.Persist(a.getEntity())
	if err == nil {
		a.id = id
	}
	return err
}
