package domain

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// LaneEntity type
type LaneEntity struct {
	ID          kernel.Id
	Name        string
	Description string
}

// LaneAggregate interface
type LaneAggregate interface {
	Saver
}

type laneAggregate struct {
	Repository
	EventRegistry
	id          kernel.Id
	name        string
	description string
}

func (a *laneAggregate) getEntity() LaneEntity {
	return LaneEntity{
		ID:          a.id,
		Name:        a.name,
		Description: a.description,
	}
}

// Save changes
func (a *laneAggregate) Save() error {
	id, err := a.Repository.Persist(a.getEntity())
	if err == nil {
		a.id = id
	}
	return err
}
