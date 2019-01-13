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
		arg2 event.Bus
		err  error
	}

	validID := kernel.ID("test")
	kind := kernel.LKind
	event.Execute(func(bus event.Bus) error {

		tests := []testcase{
			{kernel.EmptyID, kind, bus, err.ErrInvalidID},
			{validID, "", bus, err.ErrInvalidArgument},
			{validID, kind, nil, err.ErrInvalidArgument},
			{validID, kind, bus, nil},
		}

		for _, c := range tests {
			_, err := lane.Create(c.arg0, c.arg1, c.arg2)
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
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

	event.Execute(func(bus event.Bus) error {

		eventsCount := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			actual, ok := event.(lane.CreatedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
		}))

		_, err := lane.Create(validID, kind, bus)
		test.Ok(t, err)
		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestCreateLaneDefaults(t *testing.T) {
	validID := kernel.ID("test")

	event.Execute(func(bus event.Bus) error {

		entity, err := lane.Create(validID, kernel.LKind, bus)
		test.Ok(t, err)

		test.AssertExpAct(t, entity.ID, validID)
		test.AssertExpAct(t, entity.Kind, kernel.LKind)
		test.AssertExpAct(t, entity.Name, "")
		test.AssertExpAct(t, entity.Description, "")
		test.AssertExpAct(t, entity.Layout, kernel.VLayout)

		return nil
	})
}

func TestNewLane(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		arg1 event.Bus
		err  error
	}

	validID := kernel.ID("test")
	event.Execute(func(bus event.Bus) error {

		tests := []testcase{
			{kernel.EmptyID, bus, err.ErrInvalidID},
			{validID, nil, err.ErrInvalidArgument},
			{validID, bus, nil},
		}

		for _, c := range tests {
			_, err := lane.New(lane.Entity{ID: c.arg0}, c.arg1)
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
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

	event.Execute(func(bus event.Bus) error {

		aggregate, err := lane.New(lane.Entity{ID: validID, Kind: kernel.LKind}, bus)
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

		return nil
	})
}

func TestDeleteLane(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		arg1 event.Bus
		err  error
	}

	validID := kernel.ID("test")
	event.Execute(func(bus event.Bus) error {

		tests := []testcase{
			{kernel.EmptyID, bus, err.ErrInvalidID},
			{validID, nil, err.ErrInvalidArgument},
			{validID, bus, nil},
		}

		for _, c := range tests {
			err := lane.Delete(lane.Entity{ID: c.arg0}, c.arg1)
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteLaneEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := lane.Entity{ID: validID}

	expected := lane.DeletedEvent{Entity: entity}

	event.Execute(func(bus event.Bus) error {

		eventsCount := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			actual, ok := event.(lane.DeletedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
		}))

		test.Ok(t, lane.Delete(entity, bus))
		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestUpdateLaneEvents(t *testing.T) {
	validID := kernel.ID("test")
	kind := kernel.LKind

	entity := lane.Entity{ID: validID, Kind: kind, Layout: kernel.VLayout}

	event.Execute(func(bus event.Bus) error {

		aggregate, err := lane.New(entity, bus)
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

		bus.Listen(event.HandleFunc(func(event interface{}) {
			test.AssertExpAct(t, events[index], event)
			test.Assert(t, index < len(events), "Fired events count is above expectation")
			index++
		}))

		aggregate.Save()

		test.AssertExpAct(t, len(events), index)

		return nil
	})
}
