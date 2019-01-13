package services_test

import (
	"encoding/json"
	"testing"

	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"

	messagemocks "github.com/dmibod/kanban/shared/message/mocks"
	"github.com/dmibod/kanban/shared/tools/test"
	"github.com/stretchr/testify/mock"

	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

func TestShouldPublishNotification(t *testing.T) {
	id := kernel.ID("test")

	publisher := &messagemocks.Publisher{}
	publisher.On("Publish", mock.Anything).Return(nil).Once()

	err := event.Execute(func(bus event.Bus) error {
		service := services.CreateNotificationService(publisher, &noop.Logger{})
		service.Listen(bus)
		
		aggregate, err := board.New(board.Entity{ID: id}, bus)
		test.Ok(t, err)

		err = aggregate.Name("Test")
		test.Ok(t, err)

		aggregate.Save()
		return nil
	})
	test.Ok(t, err)

	publisher.AssertExpectations(t)
}

func TestShouldCollapseNotifications(t *testing.T) {
	id := kernel.ID("test")

	notifications := []kernel.Notification{kernel.Notification{Context: id, ID: id, Type: kernel.RefreshBoardNotification}}
	expected, err := json.Marshal(notifications)
	test.Ok(t, err)

	publisher := &messagemocks.Publisher{}
	publisher.On("Publish", expected).Return(nil).Once()

	err = event.Execute(func(bus event.Bus) error {
		service := services.CreateNotificationService(publisher, &noop.Logger{})
		service.Listen(bus)

		aggregate, err := board.New(board.Entity{ID: id}, bus)
		test.Ok(t, err)
		test.Ok(t, aggregate.Name("Test"))

		err = aggregate.Name("Test")
		test.Ok(t, err)

		aggregate.Save()
		return nil		
	})
	test.Ok(t, err)

	publisher.AssertExpectations(t)
}
