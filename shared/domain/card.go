package domain

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// CardCreatedEvent type
type CardCreatedEvent struct {
	ID kernel.ID
}

// CardDeletedEvent type
type CardDeletedEvent struct {
	ID kernel.ID
}

// CardNameChangedEvent type
type CardNameChangedEvent struct {
	ID       kernel.ID
	OldValue string
	NewValue string
}

// CardDescriptionChangedEvent type
type CardDescriptionChangedEvent struct {
	ID       kernel.ID
	OldValue string
	NewValue string
}

// CardEntity type
type CardEntity struct {
	ID          kernel.ID
	Name        string
	Description string
}

// Card interface
type Card interface {
	GetID() kernel.ID
	GetName() string
	GetDescription() string
	Name(string) error
	Description(string) error
}

// CardAggregate interface
type CardAggregate interface {
	Card
	Saver
}

type cardAggregate struct {
	Repository
	EventRegistry
	id          kernel.ID
	name        string
	description string
}

// NewCard aggregate
func NewCard(r Repository, e EventRegistry) (CardAggregate, error) {
	if r == nil || e == nil {
		return nil, ErrInvalidArgument
	}

	return &cardAggregate{
		Repository:    r,
		EventRegistry: e,
	}, nil
}

// DeleteCard aggregate
func DeleteCard(id kernel.ID, r Repository, e EventRegistry) (*CardEntity, error) {
	if !id.IsValid() {
		return nil, ErrInvalidID
	}

	if r == nil || e == nil {
		return nil, ErrInvalidArgument
	}

	entity, err := r.Fetch(id)
	if err != nil {
		return nil, err
	}

	card, ok := entity.(*CardEntity)
	if !ok {
		return nil, ErrInvalidType
	}

	e.Register(CardDeletedEvent{id})

	return card, nil
}

// LoadCard aggregate
func LoadCard(id kernel.ID, r Repository, e EventRegistry) (CardAggregate, error) {
	if !id.IsValid() {
		return nil, ErrInvalidID
	}

	if r == nil || e == nil {
		return nil, ErrInvalidArgument
	}

	entity, err := r.Fetch(id)
	if err != nil {
		return nil, err
	}

	card, ok := entity.(*CardEntity)
	if !ok {
		return nil, ErrInvalidType
	}

	aggregate := &cardAggregate{
		Repository:    r,
		EventRegistry: e,
	}

	aggregate.entity(*card)

	return aggregate, nil
}

// GetID
func (a *cardAggregate) GetID() kernel.ID {
	return a.id
}

// GetName
func (a *cardAggregate) GetName() string {
	return a.name
}

// Name update
func (a *cardAggregate) Name(value string) error {
	if a.name == value {
		return nil
	}

	event := CardNameChangedEvent{
		ID:       a.id,
		OldValue: a.name,
		NewValue: value,
	}
	a.register(event)

	a.name = value

	return nil
}

// GetDescription
func (a *cardAggregate) GetDescription() string {
	return a.description
}

// Description update
func (a *cardAggregate) Description(value string) error {
	if a.description == value {
		return nil
	}

	event := CardDescriptionChangedEvent{
		ID:       a.id,
		OldValue: a.description,
		NewValue: value,
	}
	a.register(event)

	a.description = value

	return nil
}

func (a *cardAggregate) getEntity() CardEntity {
	return CardEntity{
		ID:          a.id,
		Name:        a.name,
		Description: a.description,
	}
}

func (a *cardAggregate) entity(e CardEntity) {
	a.id = e.ID
	a.name = e.Name
	a.description = e.Description
}

func (a *cardAggregate) register(event interface{}) {
	if a.id.IsValid() {
		a.Register(event)
	}
}

// Save changes
func (a *cardAggregate) Save() error {
	id, err := a.Repository.Persist(a.getEntity())
	if err == nil {
		if !a.id.IsValid() {
			a.Register(CardCreatedEvent{id})
		}
		a.id = id
	}
	return err
}
