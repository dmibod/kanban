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

	service := services.CreateNotificationService(publisher, &noop.Logger{})

	err := service.Execute(func(registry *event.Manager) error {
		aggregate, err := board.New(board.Entity{ID: id}, registry)
		test.Ok(t, err)
		return aggregate.Name("Test")
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

	service := services.CreateNotificationService(publisher, &noop.Logger{})

	err = service.Execute(func(registry *event.Manager) error {
		aggregate, err := board.New(board.Entity{ID: id}, registry)
		test.Ok(t, err)
		test.Ok(t, aggregate.Name("Test"))
		return aggregate.Name("Test")
	})
	test.Ok(t, err)

	publisher.AssertExpectations(t)
}
