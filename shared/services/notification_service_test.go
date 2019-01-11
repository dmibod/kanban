package services_test

import (
	"encoding/json"
	"github.com/dmibod/kanban/shared/kernel"
	"testing"

	"github.com/dmibod/kanban/shared/domain"
	domainmocks "github.com/dmibod/kanban/shared/domain/mocks"
	messagemocks "github.com/dmibod/kanban/shared/message/mocks"
	"github.com/dmibod/kanban/shared/tools/test"
	"github.com/stretchr/testify/mock"

	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

func TestShouldPublishNotification(t *testing.T) {
	publisher := &messagemocks.Publisher{}
	publisher.On("Publish", mock.Anything).Return(nil).Once()

	repository := &domainmocks.Repository{}
	repository.On("Fetch", mock.Anything).Return(&domain.BoardEntity{}, nil)

	service := services.CreateNotificationService(publisher, &noop.Logger{})

	err := service.Execute(func(registry domain.EventRegistry) error {
		aggregate, err := domain.LoadBoard(kernel.ID("test"), repository, registry)
		test.Ok(t, err)
		return aggregate.Name("Test")
	})
	test.Ok(t, err)

	publisher.AssertExpectations(t)
}

func TestShouldCollapseNotifications(t *testing.T) {
	id := kernel.ID("test")

	notifications := []kernel.Notification{	kernel.Notification{Context:id,ID:id,Type: kernel.RefreshBoardNotification} }
	expected, err := json.Marshal(notifications)
	test.Ok(t, err)

	publisher := &messagemocks.Publisher{}
	publisher.On("Publish", expected).Return(nil).Once()

	repository := &domainmocks.Repository{}
	repository.On("Fetch", mock.Anything).Return(&domain.BoardEntity{ID:id}, nil)

	service := services.CreateNotificationService(publisher, &noop.Logger{})

	err = service.Execute(func(registry domain.EventRegistry) error {
		aggregate, err := domain.LoadBoard(id, repository, registry)
		test.Ok(t, err)
		test.Ok(t, aggregate.Name("Test"))
		return aggregate.Name("Test")
	})
	test.Ok(t, err)

	publisher.AssertExpectations(t)
}
