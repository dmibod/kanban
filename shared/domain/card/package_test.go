package card_test

import (
	"testing"

	"github.com/dmibod/kanban/shared/domain/card"
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/test"
)

func TestCreateCard(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		arg1 event.Registry
		err  error
	}

	validID := kernel.ID("test")
	manager := event.CreateEventManager()

	tests := []testcase{
		{kernel.EmptyID, manager, err.ErrInvalidID},
		{validID, manager, nil},
	}

	for _, c := range tests {
		_, err := card.Create(c.arg0, c.arg1)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestCreateCardEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID}

	expected := card.CreatedEvent{Entity: entity}

	manager := event.CreateEventManager()

	eventsCount := 0

	manager.Listen(event.HandleFunc(func(event interface{}) {
		actual, ok := event.(card.CreatedEvent)
		test.Assert(t, ok, "invalid type")
		test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
		eventsCount++
	}))

	_, err := card.Create(validID, manager)
	test.Ok(t, err)

	manager.Fire()

	test.AssertExpAct(t, 1, eventsCount)
}

func TestCreateCardDefaults(t *testing.T) {
	validID := kernel.ID("test")

	entity, err := card.Create(validID, event.CreateEventManager())
	test.Ok(t, err)

	test.AssertExpAct(t, entity.ID, validID)
	test.AssertExpAct(t, entity.Name, "")
	test.AssertExpAct(t, entity.Description, "")
}

func TestNewCard(t *testing.T) {

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
		_, err := card.New(card.Entity{ID: c.arg0}, c.arg1)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestDeleteCard(t *testing.T) {

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
		err := card.Delete(card.Entity{ID: c.arg0}, c.arg1)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestDeleteCardEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID}

	expected := card.DeletedEvent{Entity: entity}

	manager := event.CreateEventManager()

	eventsCount := 0

	manager.Listen(event.HandleFunc(func(event interface{}) {
		actual, ok := event.(card.DeletedEvent)
		test.Assert(t, ok, "invalid type")
		test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
		eventsCount++
	}))

	test.Ok(t, card.Delete(entity, manager))

	manager.Fire()

	test.AssertExpAct(t, 1, eventsCount)
}

func TestUpdateCardEvents(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID}

	manager := event.CreateEventManager()

	aggregate, err := card.New(entity, manager)
	test.Ok(t, err)

	test.Ok(t, aggregate.Name(""))
	test.Ok(t, aggregate.Name("Test"))
	test.Ok(t, aggregate.Name("Test"))

	test.Ok(t, aggregate.Description(""))
	test.Ok(t, aggregate.Description("Test"))
	test.Ok(t, aggregate.Description("Test"))

	events := []interface{}{
		card.NameChangedEvent{
			ID:       validID,
			OldValue: "",
			NewValue: "Test",
		},
		card.DescriptionChangedEvent{
			ID:       validID,
			OldValue: "",
			NewValue: "Test",
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
