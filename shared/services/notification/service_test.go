package notification_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
	messagemocks "github.com/dmibod/kanban/shared/message/mocks"
	"github.com/dmibod/kanban/shared/tools/test"
	"github.com/stretchr/testify/mock"

	"github.com/dmibod/kanban/shared/services/notification"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

func TestShouldPublishNotification(t *testing.T) {
	id := kernel.ID("test")

	publisher := &messagemocks.Publisher{}
	publisher.On("Publish", mock.Anything).Return(nil).Once()

	test.Ok(t, event.Execute(func(bus event.Bus) error {
		service := notification.CreateService(publisher, &noop.Logger{})
		service.Listen(bus)

		domainService := board.CreateService(bus)

		aggregate, err := domainService.Get(board.Entity{ID: id})

		test.Ok(t, err)
		test.Ok(t, aggregate.Name("Test"))

		bus.Fire(context.TODO())

		return nil
	}))

	publisher.AssertExpectations(t)
}

func TestShouldCollapseNotifications(t *testing.T) {
	id := kernel.ID("test")

	notifications := []kernel.Notification{kernel.Notification{BoardID: id, ID: id, Type: kernel.RefreshBoardNotification}}
	expected, err := json.Marshal(notifications)
	test.Ok(t, err)

	publisher := &messagemocks.Publisher{}
	publisher.On("Publish", expected).Return(nil).Once()

	test.Ok(t, event.Execute(func(bus event.Bus) error {
		service := notification.CreateService(publisher, &noop.Logger{})
		service.Listen(bus)

		aggregate, err := board.CreateService(bus).Get(board.Entity{ID: id})

		test.Ok(t, err)
		test.Ok(t, aggregate.Name("Test"))
		test.Ok(t, aggregate.Name("Test"))

		bus.Fire(context.TODO())

		return nil
	}))

	publisher.AssertExpectations(t)
}
