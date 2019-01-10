package domain

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// BoardNameChangedEvent type
type BoardNameChangedEvent struct {
	ID       kernel.Id
	OldValue string
	NewValue string
}

// BoardDescriptionChangedEvent type
type BoardDescriptionChangedEvent struct {
	ID       kernel.Id
	OldValue string
	NewValue string
}

// BoardLayoutChangedEvent type
type BoardLayoutChangedEvent struct {
	ID       kernel.Id
	OldValue string
	NewValue string
}

// BoardSharedChangedEvent type
type BoardSharedChangedEvent struct {
	ID       kernel.Id
	OldValue bool
	NewValue bool
}

// BoardChildAppendedEvent type
type BoardChildAppendedEvent struct {
	ID      kernel.Id
	ChildID kernel.Id
}

// BoardChildRemovedEvent type
type BoardChildRemovedEvent struct {
	ID      kernel.Id
	ChildID kernel.Id
}

// BoardEntity type
type BoardEntity struct {
	ID          kernel.Id
	Owner       string
	Name        string
	Description string
	Layout      string
	Shared      bool
	Children    []kernel.Id
}

// Board interface
type Board interface {
	GetID() kernel.Id
	GetOwner() string
	GetName() string
	Name(string) error
	GetDescription() string
	Description(string) error
	GetLayout() string
	Layout(string) error
	IsShared() bool
	Shared(bool) error
	AppendChild(kernel.Id) error
	RemoveChild(kernel.Id) error
}

// BoardAggregate interface
type BoardAggregate interface {
	Board
	Saver
}

type boardAggregate struct {
	Repository
	EventRegistry
	id          kernel.Id
	owner       string
	name        string
	description string
	layout      string
	shared      bool
	children    []kernel.Id
}

// NewBoard aggregate
func NewBoard(r Repository, e EventRegistry) (BoardAggregate, error) {
	if r == nil || e == nil {
		return nil, ErrInvalidArgument
	}

	return &boardAggregate{
		Repository:    r,
		EventRegistry: e,
	}, nil
}

// LoadBoard aggregate
func LoadBoard(id kernel.Id, r Repository, e EventRegistry) (BoardAggregate, error) {
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

	board, ok := entity.(*BoardEntity)
	if !ok {
		return nil, ErrInvalidType
	}

	aggregate := &boardAggregate{
		Repository:    r,
		EventRegistry: e,
	}

	aggregate.entity(board)

	return aggregate, nil
}

// GetID
func (a *boardAggregate) GetID() kernel.Id {
	return a.id
}

// GetOwner
func (a *boardAggregate) GetOwner() string {
	return a.owner
}

// GetName
func (a *boardAggregate) GetName() string {
	return a.name
}

// Name update
func (a *boardAggregate) Name(value string) error {
	if a.name == value {
		return nil
	}

	event := &BoardNameChangedEvent{
		ID:       a.id,
		OldValue: a.name,
		NewValue: value,
	}
	a.Register(event)

	a.name = value

	return nil
}

// GetDescription
func (a *boardAggregate) GetDescription() string {
	return a.description
}

// Description update
func (a *boardAggregate) Description(value string) error {
	if a.description == value {
		return nil
	}

	event := &BoardDescriptionChangedEvent{
		ID:       a.id,
		OldValue: a.description,
		NewValue: value,
	}
	a.Register(event)

	a.description = value

	return nil
}

// GetLayout
func (a *boardAggregate) GetLayout() string {
	return a.layout
}

// Layout update
func (a *boardAggregate) Layout(value string) error {
	if a.layout == value {
		return nil
	}

	if value == kernel.VLayout || value == kernel.HLayout {
		event := &BoardLayoutChangedEvent{
			ID:       a.id,
			OldValue: a.layout,
			NewValue: value,
		}
		a.Register(event)
		a.layout = value
	}

	return ErrInvalidLayout
}

// IsShared
func (a *boardAggregate) IsShared() bool {
	return a.shared
}

// Shared update
func (a *boardAggregate) Shared(value bool) error {
	if a.shared == value {
		return nil
	}

	event := &BoardSharedChangedEvent{
		ID:       a.id,
		OldValue: a.shared,
		NewValue: value,
	}
	a.Register(event)
	a.shared = value

	return nil
}

// AppendChild to board
func (a *boardAggregate) AppendChild(id kernel.Id) error {
	if !id.IsValid() {
		return ErrInvalidID
	}

	i := a.findChild(id)
	if i < 0 {
		a.children = append(a.children, id)

		event := &BoardChildAppendedEvent{
			ID:      a.id,
			ChildID: id,
		}
		a.Register(event)
	}

	return nil
}

// RemoveChild to board
func (a *boardAggregate) RemoveChild(id kernel.Id) error {
	if !id.IsValid() {
		return ErrInvalidID
	}

	i := a.findChild(id)
	if i < 0 {
		event := &BoardChildRemovedEvent{
			ID:      a.id,
			ChildID: a.children[i],
		}
		a.Register(event)

		a.children = append(a.children[:i], a.children[i+1:]...)
	}

	return nil
}

func (a *boardAggregate) findChild(id kernel.Id) int {
	for i, childID := range a.children {
		if childID == id {
			return i
		}
	}
	return -1
}

func (a *boardAggregate) getEntity() *BoardEntity {
	children := append([]kernel.Id{}, a.children...)
	return &BoardEntity{
		ID:          a.id,
		Owner:       a.owner,
		Name:        a.name,
		Description: a.description,
		Layout:      a.layout,
		Shared:      a.shared,
		Children:    children,
	}
}

func (a *boardAggregate) entity(e *BoardEntity) {
	if e == nil {
		return
	}

	a.id = e.ID
	a.owner = e.Owner
	a.name = e.Name
	a.description = e.Description
	a.layout = e.Layout
	a.shared = e.Shared
	a.children = append([]kernel.Id{}, e.Children...)
}

// Save changes
func (a *boardAggregate) Save() error {
	id, err := a.Repository.Persist(a.getEntity())
	if err == nil {
		a.id = id
	}
	return err
}
