package lane_test

import (
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/domain/lane"
	mocks "github.com/dmibod/kanban/shared/domain/lane/mocks"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/test"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateLane(t *testing.T) {

	type testcase struct {
		arg0 kernel.ID
		arg1 string
		err  error
	}

	validID := kernel.ID("test")
	kind := kernel.LKind

	repository := &mocks.Repository{}
	repository.On("Create", mock.Anything).Return(nil)

	event.Execute(func(bus event.Bus) error {

		tests := []testcase{
			{kernel.EmptyID, kind, err.ErrInvalidID},
			{validID, "", err.ErrInvalidArgument},
			{validID, kind, nil},
		}

		for _, c := range tests {
			domainService := lane.CreateService(repository, bus)

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
		ID:       validID,
		Kind:     kind,
		Layout:   kernel.VLayout,
		Children: []kernel.ID{},
	}

	repository := &mocks.Repository{}
	repository.On("Create", mock.Anything).Return(nil).Once()

	expected := lane.CreatedEvent{Entity: entity}

	event.Execute(func(bus event.Bus) error {

		eventsCount := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			actual, ok := event.(lane.CreatedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
		}))

		domainService := lane.CreateService(repository, bus)

		_, err := domainService.Create(validID, kind)
		test.Ok(t, err)
		test.AssertExpAct(t, 1, eventsCount)

		return nil
	})
}

func TestCreateLaneDefaults(t *testing.T) {
	validID := kernel.ID("test")

	repository := &mocks.Repository{}
	repository.On("Create", mock.Anything).Return(nil).Once()

	event.Execute(func(bus event.Bus) error {

		domainService := lane.CreateService(repository, bus)

		entity, err := domainService.Create(validID, kernel.LKind)
		test.Ok(t, err)

		test.AssertExpAct(t, entity.ID, validID)
		test.AssertExpAct(t, entity.Kind, kernel.LKind)
		test.AssertExpAct(t, entity.Name, "")
		test.AssertExpAct(t, entity.Description, "")
		test.AssertExpAct(t, entity.Layout, kernel.VLayout)

		return nil
	})
}

func TestGetLane(t *testing.T) {

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
			domainService := lane.CreateService(repository, bus)

			_, err := domainService.Get(lane.Entity{ID: c.arg0})
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteLane(t *testing.T) {

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
			domainService := lane.CreateService(repository, bus)

			err := domainService.Delete(lane.Entity{ID: c.arg0})
			test.AssertExpAct(t, c.err, err)
		}

		return nil
	})
}

func TestDeleteLaneEvent(t *testing.T) {
	validID := kernel.ID("test")
	entity := lane.Entity{ID: validID}

	expected := lane.DeletedEvent{Entity: entity}

	repository := &mocks.Repository{}
	repository.On("Delete", mock.Anything).Return(nil).Once()

	event.Execute(func(bus event.Bus) error {

		eventsCount := 0

		bus.Listen(event.HandleFunc(func(event interface{}) {
			actual, ok := event.(lane.DeletedEvent)
			test.Assert(t, ok, "invalid type")
			test.AssertExpAct(t, expected.Entity.ID, actual.Entity.ID)
			eventsCount++
		}))

		domainService := lane.CreateService(repository, bus)

		test.Ok(t, domainService.Delete(entity))
		test.AssertExpAct(t, 1, eventsCount)

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

	repository := &mocks.Repository{}

	event.Execute(func(bus event.Bus) error {
		domainService := lane.CreateService(repository, bus)

		aggregate, err := domainService.Get(lane.Entity{ID: validID, Kind: kernel.LKind})
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

	entity := lane.Entity{ID: validID, Kind: kind, Layout: kernel.VLayout}

	repository := &mocks.Repository{}
	repository.On("Update", mock.Anything).Return(nil)

	event.Execute(func(bus event.Bus) error {
		domainService := lane.CreateService(repository, bus)

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

		test.Ok(t, aggregate.Save())

		test.AssertExpAct(t, len(events), index)

		return nil
	})
}
