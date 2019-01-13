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
		arg2 event.Registry
		err  error
	}

	validID := kernel.ID("test")
	owner := "test"
	manager := event.CreateEventManager()

	tests := []testcase{
		{kernel.EmptyID, owner, manager, err.ErrInvalidID},
		{validID, "", manager, err.ErrInvalidArgument},
		{validID, owner, nil, err.ErrInvalidArgument},
		{validID, owner, manager, nil},
	}

	for _, c := range tests {
		_, err := board.Create(c.arg0, c.arg1, c.arg2)
		test.AssertExpAct(t, c.err, err)
	}
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

	manager := event.CreateEventManager()

	eventsCount := 0

	manager.Listen(event.HandleFunc(func(event interface{}) {
		actual, ok := event.(board.CreatedEvent)
		test.Assert(t, ok, "invalid type")
		test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
		eventsCount++
	}))

	_, err := board.Create(validID, owner, manager)
	test.Ok(t, err)

	manager.Fire()

	test.AssertExpAct(t, 1, eventsCount)
}

func TestCreateBoardDefaults(t *testing.T) {
	validID := kernel.ID("test")

	entity, err := board.Create(validID, "test", event.CreateEventManager())
	test.Ok(t, err)

	test.AssertExpAct(t, entity.ID, validID)
	test.AssertExpAct(t, entity.Owner, "test")
	test.AssertExpAct(t, entity.Name, "")
	test.AssertExpAct(t, entity.Description, "")
	test.AssertExpAct(t, entity.Shared, false)
	test.AssertExpAct(t, entity.Layout, kernel.VLayout)
}

func TestNewBoard(t *testing.T) {

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
		_, err := board.New(board.Entity{ID: c.arg0}, c.arg1)
		test.AssertExpAct(t, c.err, err)
	}
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

	aggregate, err := board.New(board.Entity{ID: validID, Owner: "test"}, event.CreateEventManager())
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
}

func TestDeleteBoard(t *testing.T) {

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
		err := board.Delete(board.Entity{ID: c.arg0}, c.arg1)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestDeleteBoardEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := board.Entity{ID: validID}

	expected := board.DeletedEvent{Entity: entity}

	manager := event.CreateEventManager()

	eventsCount := 0

	manager.Listen(event.HandleFunc(func(event interface{}) {
		actual, ok := event.(board.DeletedEvent)
		test.Assert(t, ok, "invalid type")
		test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
		eventsCount++
	}))

	test.Ok(t, board.Delete(entity, manager))

	manager.Fire()

	test.AssertExpAct(t, 1, eventsCount)
}

func TestUpdateBoardEvents(t *testing.T) {
	validID := kernel.ID("test")
	owner := "test"

	entity := board.Entity{ID: validID, Owner: owner, Layout: kernel.VLayout}

	manager := event.CreateEventManager()

	aggregate, err := board.New(entity, manager)
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

	manager.Listen(event.HandleFunc(func(event interface{}) {
		test.AssertExpAct(t, events[index], event)
		test.Assert(t, index < len(events), "Fired events count is above expectation")
		index++
	}))

	manager.Fire()

	test.AssertExpAct(t, len(events), index)
}
