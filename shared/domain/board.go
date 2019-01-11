package domain

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// BoardCreatedEvent type
type BoardCreatedEvent struct {
	ID kernel.ID
}

// BoardDeletedEvent type
type BoardDeletedEvent struct {
	ID kernel.ID
}

// BoardNameChangedEvent type
type BoardNameChangedEvent struct {
	ID       kernel.ID
	OldValue string
	NewValue string
}

// BoardDescriptionChangedEvent type
type BoardDescriptionChangedEvent struct {
	ID       kernel.ID
	OldValue string
	NewValue string
}

// BoardLayoutChangedEvent type
type BoardLayoutChangedEvent struct {
	ID       kernel.ID
	OldValue string
	NewValue string
}

// BoardSharedChangedEvent type
type BoardSharedChangedEvent struct {
	ID       kernel.ID
	OldValue bool
	NewValue bool
}

// BoardChildAppendedEvent type
type BoardChildAppendedEvent struct {
	ID      kernel.ID
	ChildID kernel.ID
}

// BoardChildRemovedEvent type
type BoardChildRemovedEvent struct {
	ID      kernel.ID
	ChildID kernel.ID
}

// BoardEntity type
type BoardEntity struct {
	ID          kernel.ID
	Owner       string
	Name        string
	Description string
	Layout      string
	Shared      bool
	Children    []kernel.ID
}

// Board interface
type Board interface {
	GetID() kernel.ID
	GetOwner() string
	GetName() string
	GetDescription() string
	GetLayout() string
	IsShared() bool
	Name(string) error
	Description(string) error
	Layout(string) error
	Shared(bool) error
	AppendChild(kernel.ID) error
	RemoveChild(kernel.ID) error
}

// BoardAggregate interface
type BoardAggregate interface {
	Board
	Saver
}

type boardAggregate struct {
	Repository
	EventRegistry
	id          kernel.ID
	owner       string
	name        string
	description string
	layout      string
	shared      bool
	children    []kernel.ID
}

// NewBoard aggregate
func NewBoard(owner string, r Repository, e EventRegistry) (BoardAggregate, error) {
	if owner == "" || r == nil || e == nil {
		return nil, ErrInvalidArgument
	}

	return &boardAggregate{
		Repository:    r,
		EventRegistry: e,
		owner:         owner,
		layout:        kernel.VLayout,
		shared:        false,
		children:      []kernel.ID{},
	}, nil
}

// DeleteBoard aggregate
func DeleteBoard(id kernel.ID, r Repository, e EventRegistry) (*BoardEntity, error) {
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

	e.Register(BoardDeletedEvent{id})

	return board, nil
}

// LoadBoard aggregate
func LoadBoard(id kernel.ID, r Repository, e EventRegistry) (BoardAggregate, error) {
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

	aggregate.entity(*board)

	return aggregate, nil
}

// GetID
func (a *boardAggregate) GetID() kernel.ID {
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

	event := BoardNameChangedEvent{
		ID:       a.id,
		OldValue: a.name,
		NewValue: value,
	}
	a.register(event)

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

	event := BoardDescriptionChangedEvent{
		ID:       a.id,
		OldValue: a.description,
		NewValue: value,
	}
	a.register(event)

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
		event := BoardLayoutChangedEvent{
			ID:       a.id,
			OldValue: a.layout,
			NewValue: value,
		}
		a.register(event)
		a.layout = value
		return nil
	}

	return ErrInvalidArgument
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

	event := BoardSharedChangedEvent{
		ID:       a.id,
		OldValue: a.shared,
		NewValue: value,
	}
	a.register(event)
	a.shared = value

	return nil
}

// AppendChild to board
func (a *boardAggregate) AppendChild(id kernel.ID) error {
	if !id.IsValid() {
		return ErrInvalidID
	}

	i := a.findChild(id)
	if i < 0 {
		a.children = append(a.children, id)

		event := BoardChildAppendedEvent{
			ID:      a.id,
			ChildID: id,
		}
		a.register(event)
	}

	return nil
}

// RemoveChild to board
func (a *boardAggregate) RemoveChild(id kernel.ID) error {
	if !id.IsValid() {
		return ErrInvalidID
	}

	i := a.findChild(id)
	if i >= 0 {
		event := BoardChildRemovedEvent{
			ID:      a.id,
			ChildID: a.children[i],
		}
		a.register(event)

		a.children = append(a.children[:i], a.children[i+1:]...)
	}

	return nil
}

func (a *boardAggregate) findChild(id kernel.ID) int {
	for i, childID := range a.children {
		if childID == id {
			return i
		}
	}
	return -1
}

func (a *boardAggregate) getEntity() BoardEntity {
	children := append([]kernel.ID{}, a.children...)
	return BoardEntity{
		ID:          a.id,
		Owner:       a.owner,
		Name:        a.name,
		Description: a.description,
		Layout:      a.layout,
		Shared:      a.shared,
		Children:    children,
	}
}

func (a *boardAggregate) entity(e BoardEntity) {
	a.id = e.ID
	a.owner = e.Owner
	a.name = e.Name
	a.description = e.Description
	a.layout = e.Layout
	a.shared = e.Shared
	a.children = append([]kernel.ID{}, e.Children...)
}

func (a *boardAggregate) register(event interface{}) {
	if a.id.IsValid() {
		a.Register(event)
	}
}

// Save changes
func (a *boardAggregate) Save() error {
	id, err := a.Repository.Persist(a.getEntity())
	if err == nil {
		if !a.id.IsValid() {
			a.Register(BoardCreatedEvent{id})
		}
		a.id = id
	}
	return err
}
