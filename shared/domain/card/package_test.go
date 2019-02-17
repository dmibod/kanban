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
		arg1 event.Bus
		err  error
	}

	validID := kernel.ID("test")
	event.Execute(func(bus event.Bus) error {
		tests := []testcase{
			{kernel.EmptyID, bus, err.ErrInvalidID},
			{validID, bus, nil},
		}

		for _, c := range tests {
			_, err := card.Create(c.arg0, c.arg1)
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestCreateCardEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID}

	expected := card.CreatedEvent{Entity: entity}

	event.Execute(func(bus event.Bus) error {
		eventsCount := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			actual, ok := event.(card.CreatedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
		}))

		_, err := card.Create(validID, bus)
		test.Ok(t, err)
		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestCreateCardDefaults(t *testing.T) {
	validID := kernel.ID("test")

	event.Execute(func(bus event.Bus) error {
		entity, err := card.Create(validID, bus)
		test.Ok(t, err)

		test.AssertExpAct(t, entity.ID, validID)
		test.AssertExpAct(t, entity.Name, "")
		test.AssertExpAct(t, entity.Description, "")

		return nil
	})
}

func TestNewCard(t *testing.T) {

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
			_, err := card.New(card.Entity{ID: c.arg0}, c.arg1)
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteCard(t *testing.T) {

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
			err := card.Delete(card.Entity{ID: c.arg0}, c.arg1)
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteCardEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID}

	expected := card.DeletedEvent{Entity: entity}
	event.Execute(func(bus event.Bus) error {
		eventsCount := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			actual, ok := event.(card.DeletedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
		}))

		test.Ok(t, card.Delete(entity, bus))
		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestUpdateCardEvents(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID}

	event.Execute(func(bus event.Bus) error {

		aggregate, err := card.New(entity, bus)
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
