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

func TestNewLane(t *testing.T) {

	type testcase struct {
		arg0 string
		arg1 domain.Repository
		arg2 domain.EventRegistry
		err  error
	}

	tests := []testcase{
		{"", nil, nil, domain.ErrInvalidArgument},
		{kernel.LKind, nil, nil, domain.ErrInvalidArgument},
		{kernel.LKind, &mocks.Repository{}, nil, domain.ErrInvalidArgument},
		{"", &mocks.Repository{}, domain.CreateEventManager(), domain.ErrInvalidArgument},
		{kernel.LKind, &mocks.Repository{}, domain.CreateEventManager(), nil},
		{kernel.CKind, &mocks.Repository{}, domain.CreateEventManager(), nil},
	}

	for _, c := range tests {
		_, err := domain.NewLane(c.arg0, c.arg1, c.arg2)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestDeleteLaneNegative(t *testing.T) {

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
	repository.On("Fetch", mock.Anything).Return(&domain.LaneEntity{ID: validID}, nil)

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
		_, err := domain.DeleteLane(c.arg0, c.arg1, c.arg2)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestDeleteLane(t *testing.T) {
	expected := &domain.LaneEntity{ID: kernel.ID("test")}

	repository := &mocks.Repository{}
	repository.On("Fetch", mock.Anything).Return(expected, nil)
	repository.On("Persist", mock.Anything).Return(expected.ID, nil)
	repository.On("Delete", mock.Anything).Return(expected, nil)

	eventManager := domain.CreateEventManager()

	aggregate, err := domain.NewLane(kernel.LKind, repository, eventManager)
	test.Ok(t, err)
	test.Ok(t, aggregate.Save())

	actual, err := domain.DeleteLane(aggregate.GetID(), repository, eventManager)
	test.Ok(t, err)
	test.AssertExpAct(t, expected, actual)
}

func TestLoadLane(t *testing.T) {

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
	repository.On("Fetch", mock.Anything).Return(&domain.LaneEntity{ID: validID}, nil)

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
		_, err := domain.LoadLane(c.arg0, c.arg1, c.arg2)
		test.AssertExpAct(t, c.err, err)
	}
}

func TestSaveLane(t *testing.T) {
	validID := kernel.ID("test")

	entity := domain.LaneEntity{
		ID:          kernel.EmptyID,
		Kind:        kernel.LKind,
		Name:        "Test",
		Description: "Test",
		Layout:      kernel.VLayout,
		Children:    []kernel.ID{validID},
	}

	repository := &mocks.Repository{}
	repository.On("Persist", entity).Return(validID, nil).Once()

	aggregate, err := domain.NewLane(kernel.LKind, repository, domain.CreateEventManager())
	test.Ok(t, err)

	test.Ok(t, aggregate.Name("Test"))
	test.Ok(t, aggregate.Description("Test"))
	test.Ok(t, aggregate.Layout(kernel.VLayout))
	test.Ok(t, aggregate.AppendChild(validID))

	test.Ok(t, aggregate.Save())

	repository.AssertExpectations(t)
}

func TestLaneDefaults(t *testing.T) {
	aggregate, err := domain.NewLane(kernel.LKind, &mocks.Repository{}, domain.CreateEventManager())
	test.Ok(t, err)

	test.AssertExpAct(t, aggregate.GetID(), kernel.EmptyID)
	test.AssertExpAct(t, aggregate.GetKind(), kernel.LKind)
	test.AssertExpAct(t, aggregate.GetName(), "")
	test.AssertExpAct(t, aggregate.GetDescription(), "")
	test.AssertExpAct(t, aggregate.GetLayout(), kernel.VLayout)

	aggregate, err = domain.NewLane(kernel.CKind, &mocks.Repository{}, domain.CreateEventManager())
	test.Ok(t, err)

	test.AssertExpAct(t, aggregate.GetID(), kernel.EmptyID)
	test.AssertExpAct(t, aggregate.GetKind(), kernel.CKind)
	test.AssertExpAct(t, aggregate.GetName(), "")
	test.AssertExpAct(t, aggregate.GetDescription(), "")
	test.AssertExpAct(t, aggregate.GetLayout(), "")
}

func TestLaneUpdate(t *testing.T) {
	validID := kernel.ID("test")

	aggregate, err := domain.NewLane(kernel.LKind, &mocks.Repository{}, domain.CreateEventManager())
	test.Ok(t, err)

	test.Ok(t, aggregate.Name(""))
	test.Ok(t, aggregate.Name("Test"))

	test.Ok(t, aggregate.Description(""))
	test.Ok(t, aggregate.Description("Test"))

	test.Ok(t, aggregate.Layout(kernel.VLayout))
	test.Ok(t, aggregate.Layout(kernel.HLayout))

	test.AssertExpAct(t, aggregate.Layout(""), domain.ErrInvalidArgument)
	test.AssertExpAct(t, aggregate.Layout("Test"), domain.ErrInvalidArgument)

	test.Ok(t, aggregate.AppendChild(validID))
	test.Ok(t, aggregate.RemoveChild(validID))

	test.AssertExpAct(t, aggregate.AppendChild(kernel.EmptyID), domain.ErrInvalidID)
	test.AssertExpAct(t, aggregate.RemoveChild(kernel.EmptyID), domain.ErrInvalidID)

	test.AssertExpAct(t, aggregate.GetID(), kernel.EmptyID)
	test.AssertExpAct(t, aggregate.GetKind(), kernel.LKind)
	test.AssertExpAct(t, aggregate.GetName(), "Test")
	test.AssertExpAct(t, aggregate.GetDescription(), "Test")
	test.AssertExpAct(t, aggregate.GetLayout(), kernel.HLayout)
}

func TestLaneUpdateEvents(t *testing.T) {
	validID := kernel.ID("test")

	entity := &domain.LaneEntity{ID: validID, Kind: kernel.LKind, Layout: kernel.VLayout}

	repository := &mocks.Repository{}
	repository.On("Fetch", mock.Anything).Return(entity, nil)

	eventManager := domain.CreateEventManager()

	aggregate, err := domain.LoadLane(validID, repository, eventManager)
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
		domain.LaneNameChangedEvent{
			ID:       validID,
			OldValue: "",
			NewValue: "Test",
		},
		domain.LaneDescriptionChangedEvent{
			ID:       validID,
			OldValue: "",
			NewValue: "Test",
		},
		domain.LaneLayoutChangedEvent{
			ID:       validID,
			OldValue: kernel.VLayout,
			NewValue: kernel.HLayout,
		},
		domain.LaneChildAppendedEvent{
			ID:      validID,
			ChildID: validID,
		},
		domain.LaneChildRemovedEvent{
			ID:      validID,
			ChildID: validID,
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

func TestLaneCreateDeleteEvents(t *testing.T) {
	validID := kernel.ID("test")

	entity := &domain.LaneEntity{ID: validID}

	repository := &mocks.Repository{}
	repository.On("Fetch", mock.Anything).Return(entity, nil)
	repository.On("Persist", mock.Anything).Return(validID, nil)
	repository.On("Delete", mock.Anything).Return(entity, nil)

	eventManager := domain.CreateEventManager()

	aggregate, err := domain.NewLane(kernel.LKind, repository, eventManager)
	test.Ok(t, err)
	test.Ok(t, aggregate.Save())

	_, err = domain.DeleteLane(aggregate.GetID(), repository, eventManager)
	test.Ok(t, err)

	expectedEvents := []interface{}{
		domain.LaneCreatedEvent{
			ID: validID,
		},
		domain.LaneDeletedEvent{
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
