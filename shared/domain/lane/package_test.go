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
		arg0 kernel.MemberID
		arg1 string
		err  error
	}

	validID := kernel.ID("test")
	kind := kernel.LKind

	event.Execute(func(bus event.Bus) error {

		tests := []testcase{
			{kernel.EmptyID.WithSet(validID), kind, err.ErrInvalidID},
			{validID.WithSet(validID), "", err.ErrInvalidArgument},
			{validID.WithSet(validID), kind, nil},
		}

		for _, c := range tests {
			domainService := lane.CreateService(bus)

			_, err := domainService.Create(c.arg0, c.arg1)

			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestCreateLaneEvent(t *testing.T) {
	validID := kernel.ID("test")
	kind := kernel.LKind
	entity := lane.Entity{
		ID:       validID.WithSet(validID),
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

		domainService := lane.CreateService(bus)

		_, err := domainService.Create(validID.WithSet(validID), kind)
		test.Ok(t, err)

		bus.Fire()

		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestCreateLaneDefaults(t *testing.T) {
	validID := kernel.ID("test")

	event.Execute(func(bus event.Bus) error {

		domainService := lane.CreateService(bus)

		entity, err := domainService.Create(validID.WithSet(validID), kernel.LKind)
		test.Ok(t, err)

		test.AssertExpAct(t, entity.ID, validID.WithSet(validID))
		test.AssertExpAct(t, entity.Kind, kernel.LKind)
		test.AssertExpAct(t, entity.Name, "")
		test.AssertExpAct(t, entity.Description, "")
		test.AssertExpAct(t, entity.Layout, kernel.VLayout)

		return nil
	})
}

func TestGetLane(t *testing.T) {

	type testcase struct {
		arg0 kernel.MemberID
		err  error
	}

	validID := kernel.ID("test")
	event.Execute(func(bus event.Bus) error {

		tests := []testcase{
			{kernel.EmptyID.WithID(validID), err.ErrInvalidID},
			{kernel.EmptyID.WithSet(validID), err.ErrInvalidID},
			{kernel.EmptyID.WithID(kernel.EmptyID), err.ErrInvalidID},
			{kernel.EmptyID.WithSet(kernel.EmptyID), err.ErrInvalidID},
			{validID.WithSet(validID), nil},
		}

		domainService := lane.CreateService(bus)
		for _, c := range tests {
			_, err := domainService.Get(lane.Entity{ID: c.arg0})
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteLane(t *testing.T) {

	type testcase struct {
		arg0 kernel.MemberID
		err  error
	}

	validID := kernel.ID("test")
	event.Execute(func(bus event.Bus) error {

		tests := []testcase{
			{kernel.EmptyID.WithSet(validID), err.ErrInvalidID},
			{validID.WithSet(validID), nil},
		}

		domainService := lane.CreateService(bus)
		for _, c := range tests {
			err := domainService.Delete(lane.Entity{ID: c.arg0})
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteLaneEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := lane.Entity{ID: validID.WithSet(validID)}

	expected := lane.DeletedEvent{Entity: entity}

	event.Execute(func(bus event.Bus) error {

		eventsCount := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			actual, ok := event.(lane.DeletedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
		}))

		domainService := lane.CreateService(bus)

		test.Ok(t, domainService.Delete(entity))

		bus.Fire()

		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestUpdateLane(t *testing.T) {
	validID := kernel.ID("test")

	expected := lane.Entity{
		ID:          validID.WithSet(validID),
		Kind:        kernel.LKind,
		Name:        "Test",
		Description: "Test",
		Layout:      kernel.VLayout,
		Children:    []kernel.ID{validID},
	}

	event.Execute(func(bus event.Bus) error {
		domainService := lane.CreateService(bus)

		aggregate, err := domainService.Get(lane.Entity{ID: validID.WithSet(validID), Kind: kernel.LKind})
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

func TestUpdateLaneEvents(t *testing.T) {
	validID := kernel.ID("test")
	kind := kernel.LKind

	entity := lane.Entity{ID: validID.WithSet(validID), Kind: kind, Layout: kernel.VLayout}

	event.Execute(func(bus event.Bus) error {
		domainService := lane.CreateService(bus)

		aggregate, err := domainService.Get(entity)
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
				ID:       validID.WithSet(validID),
				OldValue: "",
				NewValue: "Test",
			},
			lane.DescriptionChangedEvent{
				ID:       validID.WithSet(validID),
				OldValue: "",
				NewValue: "Test",
			},
			lane.LayoutChangedEvent{
				ID:       validID.WithSet(validID),
				OldValue: kernel.VLayout,
				NewValue: kernel.HLayout,
			},
			lane.ChildAppendedEvent{
				ID:      validID.WithSet(validID),
				ChildID: validID,
			},
			lane.ChildRemovedEvent{
				ID:      validID.WithSet(validID),
				ChildID: validID,
			},
		}

		index := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			test.AssertExpAct(t, events[index], event)
			test.Assert(t, index < len(events), "Fired events count is above expectation")
			index++
		}))

		bus.Fire()

		test.AssertExpAct(t, len(events), index)

		return nil
	})
}
