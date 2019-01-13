package board_test

import (
	"testing"

	"github.com/dmibod/kanban/shared/domain/board"
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/test"
)

func TestCreateBoard(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		arg1 string
		arg2 event.Bus
		err  error
	}

	validID := kernel.ID("test")
	owner := "test"

	event.Execute(func(bus event.Bus) error {

		tests := []testcase{
			{kernel.EmptyID, owner, bus, err.ErrInvalidID},
			{validID, "", bus, err.ErrInvalidArgument},
			{validID, owner, nil, err.ErrInvalidArgument},
			{validID, owner, bus, nil},
		}

		for _, c := range tests {
			_, err := board.Create(c.arg0, c.arg1, c.arg2)
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestCreateBoardEvent(t *testing.T) {
	validID := kernel.ID("test")
	owner := "test"
	entity := board.Entity{
		ID:       validID,
		Owner:    owner,
		Layout:   kernel.VLayout,
		Shared:   false,
		Children: []kernel.ID{},
	}

	expected := board.CreatedEvent{Entity: entity}

	event.Execute(func(bus event.Bus) error {

		eventsCount := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			actual, ok := event.(board.CreatedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
		}))

		_, err := board.Create(validID, owner, bus)
		test.Ok(t, err)
		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestCreateBoardDefaults(t *testing.T) {
	validID := kernel.ID("test")

	event.Execute(func(bus event.Bus) error {

		entity, err := board.Create(validID, "test", bus)
		test.Ok(t, err)

		test.AssertExpAct(t, entity.ID, validID)
		test.AssertExpAct(t, entity.Owner, "test")
		test.AssertExpAct(t, entity.Name, "")
		test.AssertExpAct(t, entity.Description, "")
		test.AssertExpAct(t, entity.Shared, false)
		test.AssertExpAct(t, entity.Layout, kernel.VLayout)

		return nil
	})
}

func TestNewBoard(t *testing.T) {

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
			_, err := board.New(board.Entity{ID: c.arg0}, c.arg1)
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestUpdateBoard(t *testing.T) {
	validID := kernel.ID("test")

	expected := board.Entity{
		ID:          validID,
		Owner:       "test",
		Name:        "Test",
		Description: "Test",
		Layout:      kernel.VLayout,
		Shared:      true,
		Children:    []kernel.ID{validID},
	}

	event.Execute(func(bus event.Bus) error {

		aggregate, err := board.New(board.Entity{ID: validID, Owner: "test"}, bus)
		test.Ok(t, err)

		test.Ok(t, aggregate.Name("Test"))
		test.Ok(t, aggregate.Description("Test"))
		test.Ok(t, aggregate.Shared(true))
		test.Ok(t, aggregate.Layout(kernel.VLayout))
		test.Ok(t, aggregate.AppendChild(validID))

		actual := aggregate.Root()

		test.AssertExpAct(t, expected.ID, actual.ID)
		test.AssertExpAct(t, expected.Owner, actual.Owner)
		test.AssertExpAct(t, expected.Name, actual.Name)
		test.AssertExpAct(t, expected.Description, actual.Description)
		test.AssertExpAct(t, expected.Shared, actual.Shared)
		test.AssertExpAct(t, expected.Layout, actual.Layout)
		test.AssertExpAct(t, len(expected.Children), len(actual.Children))

		return nil
	})
}

func TestDeleteBoard(t *testing.T) {

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
			err := board.Delete(board.Entity{ID: c.arg0}, c.arg1)
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteBoardEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := board.Entity{ID: validID}

	expected := board.DeletedEvent{Entity: entity}

	event.Execute(func(bus event.Bus) error {

		eventsCount := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			actual, ok := event.(board.DeletedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
		}))

		test.Ok(t, board.Delete(entity, bus))
		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestUpdateBoardEvents(t *testing.T) {
	validID := kernel.ID("test")
	owner := "test"

	entity := board.Entity{ID: validID, Owner: owner, Layout: kernel.VLayout}

	event.Execute(func(bus event.Bus) error {

		aggregate, err := board.New(entity, bus)
		test.Ok(t, err)

		test.Ok(t, aggregate.Name(""))
		test.Ok(t, aggregate.Name("Test"))
		test.Ok(t, aggregate.Name("Test"))

		test.Ok(t, aggregate.Description(""))
		test.Ok(t, aggregate.Description("Test"))
		test.Ok(t, aggregate.Description("Test"))

		test.Ok(t, aggregate.Shared(false))
		test.Ok(t, aggregate.Shared(true))
		test.Ok(t, aggregate.Shared(true))

		test.Ok(t, aggregate.Layout(kernel.VLayout))
		test.Ok(t, aggregate.Layout(kernel.VLayout))

		test.Ok(t, aggregate.Layout(kernel.HLayout))
		test.Ok(t, aggregate.Layout(kernel.HLayout))

		test.Ok(t, aggregate.AppendChild(validID))
		test.Ok(t, aggregate.AppendChild(validID))

		test.Ok(t, aggregate.RemoveChild(validID))
		test.Ok(t, aggregate.RemoveChild(validID))

		events := []interface{}{
			board.NameChangedEvent{
				ID:       validID,
				OldValue: "",
				NewValue: "Test",
			},
			board.DescriptionChangedEvent{
				ID:       validID,
				OldValue: "",
				NewValue: "Test",
			},
			board.SharedChangedEvent{
				ID:       validID,
				OldValue: false,
				NewValue: true,
			},
			board.LayoutChangedEvent{
				ID:       validID,
				OldValue: kernel.VLayout,
				NewValue: kernel.HLayout,
			},
			board.ChildAppendedEvent{
				ID:      validID,
				ChildID: validID,
			},
			board.ChildRemovedEvent{
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
