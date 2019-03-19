package card_test

import (
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/dmibod/kanban/shared/domain/card"
	mocks "github.com/dmibod/kanban/shared/domain/card/mocks"
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/test"
)

func TestCreateCard(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		err  error
	}

	validID := kernel.ID("test")

	repository := &mocks.Repository{}
	repository.On("Create", mock.Anything).Return(nil)

	event.Execute(func(bus event.Bus) error {
		tests := []testcase{
			{kernel.EmptyID, err.ErrInvalidID},
			{validID, nil},
		}

		for _, c := range tests {
			domainService := card.CreateService(repository, bus)

			_, err := domainService.Create(c.arg0)
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestCreateCardEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID}

	repository := &mocks.Repository{}
	repository.On("Create", mock.Anything).Return(nil).Once()

	expected := card.CreatedEvent{Entity: entity}

	event.Execute(func(bus event.Bus) error {
		eventsCount := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			actual, ok := event.(card.CreatedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
		}))

		domainService := card.CreateService(repository, bus)

		_, err := domainService.Create(validID)
		test.Ok(t, err)
		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestCreateCardDefaults(t *testing.T) {
	validID := kernel.ID("test")

	repository := &mocks.Repository{}
	repository.On("Create", mock.Anything).Return(nil).Once()

	event.Execute(func(bus event.Bus) error {
		domainService := card.CreateService(repository, bus)

		entity, err := domainService.Create(validID)
		test.Ok(t, err)

		test.AssertExpAct(t, entity.ID, validID)
		test.AssertExpAct(t, entity.Name, "")
		test.AssertExpAct(t, entity.Description, "")

		return nil
	})
}

func TestGetCard(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		err  error
	}

	repository := &mocks.Repository{}

	validID := kernel.ID("test")
	event.Execute(func(bus event.Bus) error {
		tests := []testcase{
			{kernel.EmptyID, err.ErrInvalidID},
			{validID, nil},
		}

		for _, c := range tests {
			domainService := card.CreateService(repository, bus)

			_, err := domainService.Get(card.Entity{ID: c.arg0})
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteCard(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		err  error
	}

	repository := &mocks.Repository{}
	repository.On("Delete", mock.Anything).Return(nil).Once()

	validID := kernel.ID("test")
	event.Execute(func(bus event.Bus) error {
		tests := []testcase{
			{kernel.EmptyID, err.ErrInvalidID},
			{validID, nil},
		}

		for _, c := range tests {
			domainService := card.CreateService(repository, bus)

			err := domainService.Delete(card.Entity{ID: c.arg0})
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteCardEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID}

	expected := card.DeletedEvent{Entity: entity}

	repository := &mocks.Repository{}
	repository.On("Delete", mock.Anything).Return(nil).Once()

	event.Execute(func(bus event.Bus) error {
		eventsCount := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			actual, ok := event.(card.DeletedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
		}))

		domainService := card.CreateService(repository, bus)

		test.Ok(t, domainService.Delete(entity))
		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestUpdateCardEvents(t *testing.T) {
	validID := kernel.ID("test")
	entity := card.Entity{ID: validID}

	repository := &mocks.Repository{}
	repository.On("Update", mock.Anything).Return(nil)

	event.Execute(func(bus event.Bus) error {
		domainService := card.CreateService(repository, bus)

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

		test.Ok(t, aggregate.Save())

		test.AssertExpAct(t, len(events), index)

		return nil
	})
}
