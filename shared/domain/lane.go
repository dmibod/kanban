package domain

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// LaneNameChangedEvent type
type LaneNameChangedEvent struct {
	ID       kernel.ID
	OldValue string
	NewValue string
}

// LaneDescriptionChangedEvent type
type LaneDescriptionChangedEvent struct {
	ID       kernel.ID
	OldValue string
	NewValue string
}

// LaneLayoutChangedEvent type
type LaneLayoutChangedEvent struct {
	ID       kernel.ID
	OldValue string
	NewValue string
}

// LaneChildAppendedEvent type
type LaneChildAppendedEvent struct {
	ID      kernel.ID
	ChildID kernel.ID
}

// LaneChildRemovedEvent type
type LaneChildRemovedEvent struct {
	ID      kernel.ID
	ChildID kernel.ID
}

// LaneEntity type
type LaneEntity struct {
	ID          kernel.ID
	Kind        string
	Name        string
	Description string
	Layout      string
	Children    []kernel.ID
}

// Lane interface
type Lane interface {
	GetID() kernel.ID
	GetKind() string
	GetName() string
	GetDescription() string
	GetLayout() string
	Name(string) error
	Description(string) error
	Layout(string) error
	AppendChild(kernel.ID) error
	RemoveChild(kernel.ID) error
}

// LaneAggregate interface
type LaneAggregate interface {
	Lane
	Saver
}

type laneAggregate struct {
	Repository
	EventRegistry
	id          kernel.ID
	kind        string
	name        string
	description string
	layout      string
	children    []kernel.ID
}

// NewLane aggregate
func NewLane(kind string, r Repository, e EventRegistry) (LaneAggregate, error) {
	if r == nil || e == nil {
		return nil, ErrInvalidArgument
	}

	if kind != kernel.LKind && kind != kernel.CKind {
		return nil, ErrInvalidArgument
	}

	return &laneAggregate{
		Repository:    r,
		EventRegistry: e,
		kind:          kind,
		layout:        kernel.VLayout,
		children:      []kernel.ID{},
	}, nil
}

// LoadLane aggregate
func LoadLane(id kernel.ID, r Repository, e EventRegistry) (LaneAggregate, error) {
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

	Lane, ok := entity.(*LaneEntity)
	if !ok {
		return nil, ErrInvalidType
	}

	aggregate := &laneAggregate{
		Repository:    r,
		EventRegistry: e,
	}

	aggregate.entity(*Lane)

	return aggregate, nil
}

// GetID
func (a *laneAggregate) GetID() kernel.ID {
	return a.id
}

// GetKind
func (a *laneAggregate) GetKind() string {
	return a.kind
}

// GetName
func (a *laneAggregate) GetName() string {
	return a.name
}

// Name update
func (a *laneAggregate) Name(value string) error {
	if a.name == value {
		return nil
	}

	event := LaneNameChangedEvent{
		ID:       a.id,
		OldValue: a.name,
		NewValue: value,
	}
	a.Register(event)

	a.name = value

	return nil
}

// GetDescription
func (a *laneAggregate) GetDescription() string {
	return a.description
}

// Description update
func (a *laneAggregate) Description(value string) error {
	if a.description == value {
		return nil
	}

	event := LaneDescriptionChangedEvent{
		ID:       a.id,
		OldValue: a.description,
		NewValue: value,
	}
	a.Register(event)

	a.description = value

	return nil
}

// GetLayout
func (a *laneAggregate) GetLayout() string {
	return a.layout
}

// Layout update
func (a *laneAggregate) Layout(value string) error {
	if a.layout == value {
		return nil
	}

	if value == kernel.VLayout || value == kernel.HLayout {
		event := LaneLayoutChangedEvent{
			ID:       a.id,
			OldValue: a.layout,
			NewValue: value,
		}
		a.Register(event)
		a.layout = value
		return nil
	}

	return ErrInvalidArgument
}

// AppendChild to Lane
func (a *laneAggregate) AppendChild(id kernel.ID) error {
	if !id.IsValid() {
		return ErrInvalidID
	}

	i := a.findChild(id)
	if i < 0 {
		a.children = append(a.children, id)

		event := LaneChildAppendedEvent{
			ID:      a.id,
			ChildID: id,
		}
		a.Register(event)
	}

	return nil
}

// RemoveChild to Lane
func (a *laneAggregate) RemoveChild(id kernel.ID) error {
	if !id.IsValid() {
		return ErrInvalidID
	}

	i := a.findChild(id)
	if i >= 0 {
		event := LaneChildRemovedEvent{
			ID:      a.id,
			ChildID: a.children[i],
		}
		a.Register(event)

		a.children = append(a.children[:i], a.children[i+1:]...)
	}

	return nil
}

func (a *laneAggregate) findChild(id kernel.ID) int {
	for i, childID := range a.children {
		if childID == id {
			return i
		}
	}
	return -1
}

func (a *laneAggregate) getEntity() LaneEntity {
	children := append([]kernel.ID{}, a.children...)
	return LaneEntity{
		ID:          a.id,
		Kind:        a.kind,
		Name:        a.name,
		Description: a.description,
		Layout:      a.layout,
		Children:    children,
	}
}

func (a *laneAggregate) entity(e LaneEntity) {
	a.id = e.ID
	a.kind = e.Kind
	a.name = e.Name
	a.description = e.Description
	a.layout = e.Layout
	a.children = append([]kernel.ID{}, e.Children...)
}

// Save changes
func (a *laneAggregate) Save() error {
	id, err := a.Repository.Persist(a.getEntity())
	if err == nil {
		a.id = id
	}
	return err
}
