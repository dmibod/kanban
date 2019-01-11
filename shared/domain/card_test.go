package domain_test

import (
	"errors"
	"testing"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/stretchr/testify/mock"

	"github.com/dmibod/kanban/shared/domain/mocks"
	"github.com/dmibod/kanban/shared/tools/test"

	"github.com/dmibod/kanban/shared/domain"
)

func TestNewCard(t *testing.T) {

	type testcase struct {
		arg0 domain.Repository
		arg1 domain.EventRegistry
		err  error
	}

	tests := []testcase{
		{nil, nil, domain.ErrInvalidArgument},
		{&mocks.Repository{}, nil, domain.ErrInvalidArgument},
		{&mocks.Repository{}, domain.CreateEventManager(), nil},
	}

	for _, c := range tests {
		_, err := domain.NewCard(c.arg0, c.arg1)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestDeleteCardNegative(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		arg1 domain.Repository
		arg2 domain.EventRegistry
		err  error
	}

	validID := kernel.ID("test")
	fetchErr := errors.New("fetch error")

	eventManager := domain.CreateEventManager()

	repository := &mocks.Repository{}
	repository.On("Fetch", mock.Anything).Return(&domain.CardEntity{ID: validID}, nil)

	fetchErrRepo := &mocks.Repository{}
	fetchErrRepo.On("Fetch", mock.Anything).Return(nil, fetchErr)

	wrongResultRepo := &mocks.Repository{}
	wrongResultRepo.On("Fetch", mock.Anything).Return(&struct{}{}, nil)

	tests := []testcase{
		{kernel.EmptyID, nil, nil, domain.ErrInvalidID},
		{validID, nil, nil, domain.ErrInvalidArgument},
		{validID, repository, nil, domain.ErrInvalidArgument},
		{validID, fetchErrRepo, eventManager, fetchErr},
		{validID, wrongResultRepo, eventManager, domain.ErrInvalidType},
		{validID, repository, eventManager, nil},
	}

	for _, c := range tests {
		_, err := domain.DeleteCard(c.arg0, c.arg1, c.arg2)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestDeleteCard(t *testing.T) {
	expected := &domain.CardEntity{ID: kernel.ID("test")}

	repository := &mocks.Repository{}
	repository.On("Fetch", mock.Anything).Return(expected, nil)
	repository.On("Persist", mock.Anything).Return(expected.ID, nil)
	repository.On("Delete", mock.Anything).Return(expected, nil)

	eventManager := domain.CreateEventManager()

	aggregate, err := domain.NewCard(repository, eventManager)
	test.Ok(t, err)
	test.Ok(t, aggregate.Save())

	actual, err := domain.DeleteCard(aggregate.GetID(), repository, eventManager)
	test.Ok(t, err)
	test.AssertExpAct(t, expected, actual)
}

func TestLoadCard(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		arg1 domain.Repository
		arg2 domain.EventRegistry
		err  error
	}

	validID := kernel.ID("test")
	fetchErr := errors.New("fetch error")

	eventManager := domain.CreateEventManager()

	repository := &mocks.Repository{}
	repository.On("Fetch", mock.Anything).Return(&domain.CardEntity{ID: validID}, nil)

	fetchErrRepo := &mocks.Repository{}
	fetchErrRepo.On("Fetch", mock.Anything).Return(nil, fetchErr)

	wrongResultRepo := &mocks.Repository{}
	wrongResultRepo.On("Fetch", mock.Anything).Return(&struct{}{}, nil)

	tests := []testcase{
		{kernel.EmptyID, nil, nil, domain.ErrInvalidID},
		{validID, nil, nil, domain.ErrInvalidArgument},
		{validID, repository, nil, domain.ErrInvalidArgument},
		{validID, fetchErrRepo, eventManager, fetchErr},
		{validID, wrongResultRepo, eventManager, domain.ErrInvalidType},
		{validID, repository, eventManager, nil},
	}

	for _, c := range tests {
		_, err := domain.LoadCard(c.arg0, c.arg1, c.arg2)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestSaveCard(t *testing.T) {
	validID := kernel.ID("test")

	entity := domain.CardEntity{
		ID:          kernel.EmptyID,
		Name:        "Test",
		Description: "Test",
	}

	repository := &mocks.Repository{}
	repository.On("Persist", entity).Return(validID, nil).Once()

	aggregate, err := domain.NewCard(repository, domain.CreateEventManager())
	test.Ok(t, err)

	test.Ok(t, aggregate.Name("Test"))
	test.Ok(t, aggregate.Description("Test"))

	test.Ok(t, aggregate.Save())

	repository.AssertExpectations(t)
}

func TestCardDefaults(t *testing.T) {
	aggregate, err := domain.NewCard(&mocks.Repository{}, domain.CreateEventManager())
	test.Ok(t, err)

	test.AssertExpAct(t, aggregate.GetID(), kernel.EmptyID)
	test.AssertExpAct(t, aggregate.GetName(), "")
	test.AssertExpAct(t, aggregate.GetDescription(), "")
}

func TestCardUpdate(t *testing.T) {
	aggregate, err := domain.NewCard(&mocks.Repository{}, domain.CreateEventManager())
	test.Ok(t, err)

	test.Ok(t, aggregate.Name(""))
	test.Ok(t, aggregate.Name("Test"))

	test.Ok(t, aggregate.Description(""))
	test.Ok(t, aggregate.Description("Test"))

	test.AssertExpAct(t, aggregate.GetID(), kernel.EmptyID)
	test.AssertExpAct(t, aggregate.GetName(), "Test")
	test.AssertExpAct(t, aggregate.GetDescription(), "Test")
}

func TestCardUpdateEvents(t *testing.T) {
	validID := kernel.ID("test")

	entity := &domain.CardEntity{ID: validID}

	repository := &mocks.Repository{}
	repository.On("Fetch", mock.Anything).Return(entity, nil)

	eventManager := domain.CreateEventManager()

	aggregate, err := domain.LoadCard(validID, repository, eventManager)
	test.Ok(t, err)

	test.Ok(t, aggregate.Name(""))
	test.Ok(t, aggregate.Name("Test"))
	test.Ok(t, aggregate.Name("Test"))

	test.Ok(t, aggregate.Description(""))
	test.Ok(t, aggregate.Description("Test"))
	test.Ok(t, aggregate.Description("Test"))

	events := []interface{}{
		domain.CardNameChangedEvent{
			ID:       validID,
			OldValue: "",
			NewValue: "Test",
		},
		domain.CardDescriptionChangedEvent{
			ID:       validID,
			OldValue: "",
			NewValue: "Test",
		},
	}

	index := 0

	eventManager.Listen(domain.HandleFunc(func(event interface{}) {
		test.AssertExpAct(t, events[index], event)
		test.Assert(t, index < len(events), "Fired events count is above expectation")
		index++
	}))

	eventManager.Fire()

	test.AssertExpAct(t, len(events), index)
}

func TestCardCreateDeleteEvents(t *testing.T) {
	validID := kernel.ID("test")

	entity := &domain.CardEntity{ID: validID}

	repository := &mocks.Repository{}
	repository.On("Fetch", mock.Anything).Return(entity, nil)
	repository.On("Persist", mock.Anything).Return(validID, nil)
	repository.On("Delete", mock.Anything).Return(entity, nil)

	eventManager := domain.CreateEventManager()

	aggregate, err := domain.NewCard(repository, eventManager)
	test.Ok(t, err)
	test.Ok(t, aggregate.Save())

	_, err = domain.DeleteCard(aggregate.GetID(), repository, eventManager)
	test.Ok(t, err)

	expectedEvents := []interface{}{
		domain.CardCreatedEvent{
			ID: validID,
		},
		domain.CardDeletedEvent{
			ID: validID,
		},
	}

	index := 0

	eventManager.Listen(domain.HandleFunc(func(event interface{}) {
		test.AssertExpAct(t, expectedEvents[index], event)
		test.Assert(t, index < len(expectedEvents), "Fired events count is above expectation")
		index++
	}))

	eventManager.Fire()

	test.AssertExpAct(t, len(expectedEvents), index)
}
