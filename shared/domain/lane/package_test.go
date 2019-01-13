package lane_test

import (
	"testing"

	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/domain/lane"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/test"
)

func TestCreateLane(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		arg1 string
		arg2 event.Registry
		err  error
	}

	validID := kernel.ID("test")
	kind := kernel.LKind
	manager := event.CreateEventManager()

	tests := []testcase{
		{kernel.EmptyID, kind, manager, err.ErrInvalidID},
		{validID, "", manager, err.ErrInvalidArgument},
		{validID, kind, nil, err.ErrInvalidArgument},
		{validID, kind, manager, nil},
	}

	for _, c := range tests {
		_, err := lane.Create(c.arg0, c.arg1, c.arg2)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestCreateLaneEvent(t *testing.T) {
	validID := kernel.ID("test")
	kind := kernel.LKind
	entity := lane.Entity{
		ID:       validID,
		Kind:     kind,
		Layout:   kernel.VLayout,
		Children: []kernel.ID{},
	}

	expected := lane.CreatedEvent{Entity: entity}

	manager := event.CreateEventManager()

	eventsCount := 0

	manager.Listen(event.HandleFunc(func(event interface{}) {
		actual, ok := event.(lane.CreatedEvent)
		test.Assert(t, ok, "invalid type")
		test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
		eventsCount++
	}))

	_, err := lane.Create(validID, kind, manager)
	test.Ok(t, err)

	manager.Fire()

	test.AssertExpAct(t, 1, eventsCount)
}

func TestCreateLaneDefaults(t *testing.T) {
	validID := kernel.ID("test")

	entity, err := lane.Create(validID, kernel.LKind, event.CreateEventManager())
	test.Ok(t, err)

	test.AssertExpAct(t, entity.ID, validID)
	test.AssertExpAct(t, entity.Kind, kernel.LKind)
	test.AssertExpAct(t, entity.Name, "")
	test.AssertExpAct(t, entity.Description, "")
	test.AssertExpAct(t, entity.Layout, kernel.VLayout)
}

func TestNewLane(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		arg1 event.Registry
		err  error
	}

	validID := kernel.ID("test")
	manager := event.CreateEventManager()

	tests := []testcase{
		{kernel.EmptyID, manager, err.ErrInvalidID},
		{validID, nil, err.ErrInvalidArgument},
		{validID, manager, nil},
	}

	for _, c := range tests {
		_, err := lane.New(lane.Entity{ID: c.arg0}, c.arg1)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestUpdateLane(t *testing.T) {
	validID := kernel.ID("test")

	expected := lane.Entity{
		ID:          validID,
		Kind:        kernel.LKind,
		Name:        "Test",
		Description: "Test",
		Layout:      kernel.VLayout,
		Children:    []kernel.ID{validID},
	}

	aggregate, err := lane.New(lane.Entity{ID: validID, Kind: kernel.LKind}, event.CreateEventManager())
	test.Ok(t, err)

	test.Ok(t, aggregate.Name("Test"))
	test.Ok(t, aggregate.Description("Test"))
	test.Ok(t, aggregate.Layout(kernel.VLayout))
	test.Ok(t, aggregate.AppendChild(validID))

	actual := aggregate.Root()

	test.AssertExpAct(t, expected.ID, actual.ID)
	test.AssertExpAct(t, expected.Kind, actual.Kind)
	test.AssertExpAct(t, expected.Name, actual.Name)
	test.AssertExpAct(t, expected.Description, actual.Description)
	test.AssertExpAct(t, expected.Layout, actual.Layout)
	test.AssertExpAct(t, len(expected.Children), len(actual.Children))
}

func TestDeleteLane(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		arg1 event.Registry
		err  error
	}

	validID := kernel.ID("test")
	manager := event.CreateEventManager()

	tests := []testcase{
		{kernel.EmptyID, manager, err.ErrInvalidID},
		{validID, nil, err.ErrInvalidArgument},
		{validID, manager, nil},
	}

	for _, c := range tests {
		err := lane.Delete(lane.Entity{ID: c.arg0}, c.arg1)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestDeleteLaneEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := lane.Entity{ID: validID}

	expected := lane.DeletedEvent{Entity: entity}

	manager := event.CreateEventManager()

	eventsCount := 0

	manager.Listen(event.HandleFunc(func(event interface{}) {
		actual, ok := event.(lane.DeletedEvent)
		test.Assert(t, ok, "invalid type")
		test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
		eventsCount++
	}))

	test.Ok(t, lane.Delete(entity, manager))

	manager.Fire()

	test.AssertExpAct(t, 1, eventsCount)
}

func TestUpdateLaneEvents(t *testing.T) {
	validID := kernel.ID("test")
	kind := kernel.LKind

	entity := lane.Entity{ID: validID, Kind: kind, Layout: kernel.VLayout}

	manager := event.CreateEventManager()

	aggregate, err := lane.New(entity, manager)
	test.Ok(t, err)

	test.Ok(t, aggregate.Name(""))
	test.Ok(t, aggregate.Name("Test"))
	test.Ok(t, aggregate.Name("Test"))

	test.Ok(t, aggregate.Description(""))
	test.Ok(t, aggregate.Description("Test"))
	test.Ok(t, aggregate.Description("Test"))

	test.Ok(t, aggregate.Layout(kernel.VLayout))
	test.Ok(t, aggregate.Layout(kernel.VLayout))

	test.Ok(t, aggregate.Layout(kernel.HLayout))
	test.Ok(t, aggregate.Layout(kernel.HLayout))

	test.Ok(t, aggregate.AppendChild(validID))
	test.Ok(t, aggregate.AppendChild(validID))

	test.Ok(t, aggregate.RemoveChild(validID))
	test.Ok(t, aggregate.RemoveChild(validID))

	events := []interface{}{
		lane.NameChangedEvent{
			ID:       validID,
			OldValue: "",
			NewValue: "Test",
		},
		lane.DescriptionChangedEvent{
			ID:       validID,
			OldValue: "",
			NewValue: "Test",
		},
		lane.LayoutChangedEvent{
			ID:       validID,
			OldValue: kernel.VLayout,
			NewValue: kernel.HLayout,
		},
		lane.ChildAppendedEvent{
			ID:      validID,
			ChildID: validID,
		},
		lane.ChildRemovedEvent{
			ID:      validID,
			ChildID: validID,
		},
	}

	index := 0

	manager.Listen(event.HandleFunc(func(event interface{}) {
		test.AssertExpAct(t, events[index], event)
		test.Assert(t, index < len(events), "Fired events count is above expectation")
		index++
	}))

	manager.Fire()

	test.AssertExpAct(t, len(events), index)
}
