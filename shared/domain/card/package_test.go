package card_test

import (
	"context"
	"testing"

	"github.com/dmibod/kanban/shared/domain/card"
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/test"
)

func TestCreateCard(t *testing.T) {

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
			{validID.WithID(validID), nil},
			{validID.WithSet(validID), nil},
		}

		domainService := card.CreateService(bus)
		for _, c := range tests {
			_, err := domainService.Create(c.arg0)
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestCreateCardEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID.WithSet(validID)}

	expected := card.CreatedEvent{Entity: entity}

	event.Execute(func(bus event.Bus) error {
		eventsCount := 0

		bus.Listen(event.HandleFunc(func(ctx context.Context, event interface{}) error {
			actual, ok := event.(card.CreatedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
			return nil
		}))

		domainService := card.CreateService(bus)

		_, err := domainService.Create(validID.WithSet(validID))
		test.Ok(t, err)

		bus.Fire(context.TODO())

		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestCreateCardDefaults(t *testing.T) {
	validID := kernel.ID("test")

	event.Execute(func(bus event.Bus) error {
		domainService := card.CreateService(bus)

		entity, err := domainService.Create(validID.WithSet(validID))
		test.Ok(t, err)

		test.AssertExpAct(t, entity.ID.ID, validID)
		test.AssertExpAct(t, entity.ID.SetID, validID)
		test.AssertExpAct(t, entity.Name, "")
		test.AssertExpAct(t, entity.Description, "")

		return nil
	})
}

func TestGetCard(t *testing.T) {

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
			{validID.WithID(validID), nil},
			{validID.WithSet(validID), nil},
		}

		domainService := card.CreateService(bus)
		for _, c := range tests {
			_, err := domainService.Get(card.Entity{ID: c.arg0})
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteCard(t *testing.T) {

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
			{validID.WithID(validID), nil},
			{validID.WithSet(validID), nil},
		}

		domainService := card.CreateService(bus)

		for _, c := range tests {
			err := domainService.Delete(card.Entity{ID: c.arg0})
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteCardEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID.WithSet(validID)}

	expected := card.DeletedEvent{Entity: entity}

	event.Execute(func(bus event.Bus) error {
		eventsCount := 0

		bus.Listen(event.HandleFunc(func(ctx context.Context, event interface{}) error {
			actual, ok := event.(card.DeletedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
			return nil
		}))

		domainService := card.CreateService(bus)

		test.Ok(t, domainService.Delete(entity))

		bus.Fire(context.TODO())

		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestUpdateCardEvents(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID.WithSet(validID)}

	event.Execute(func(bus event.Bus) error {
		domainService := card.CreateService(bus)

		aggregate, err := domainService.Get(entity)
		test.Ok(t, err)

		test.Ok(t, aggregate.Name(""))
		test.Ok(t, aggregate.Name("Test"))
		test.Ok(t, aggregate.Name("Test"))

		test.Ok(t, aggregate.Description(""))
		test.Ok(t, aggregate.Description("Test"))
		test.Ok(t, aggregate.Description("Test"))

		events := []interface{}{
			card.NameChangedEvent{
				ID:       validID.WithSet(validID),
				OldValue: "",
				NewValue: "Test",
			},
			card.DescriptionChangedEvent{
				ID:       validID.WithSet(validID),
				OldValue: "",
				NewValue: "Test",
			},
		}

		index := 0

		bus.Listen(event.HandleFunc(func(ctx context.Context, event interface{}) error {
			test.AssertExpAct(t, events[index], event)
			test.Assert(t, index < len(events), "Fired events count is above expectation")
			index++
			return nil
		}))

		bus.Fire(context.TODO())

		test.AssertExpAct(t, len(events), index)

		return nil
	})
}
