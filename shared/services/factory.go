package services

import (
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services/board"
	"github.com/dmibod/kanban/shared/services/card"
	"github.com/dmibod/kanban/shared/services/command"
	"github.com/dmibod/kanban/shared/services/lane"
	"github.com/dmibod/kanban/shared/services/notification"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// ServiceFactory creates service instances
type ServiceFactory struct {
	persistence.RepositoryFactory
	message.Publisher
	logger.Logger
}

// CreateServiceFactory creates service factory
func CreateServiceFactory(f persistence.RepositoryFactory, p message.Publisher, l logger.Logger) *ServiceFactory {
	return &ServiceFactory{
		RepositoryFactory: f,
		Publisher:         p,
		Logger:            l,
	}
}

// CreateBoardService creates new service instance
func (f *ServiceFactory) CreateBoardService() board.Service {
	return board.CreateService(
		f.CreateNotificationService(),
		f.RepositoryFactory.CreateRepository(),
		f.Logger)
}

// CreateLaneService creates new service instance
func (f *ServiceFactory) CreateLaneService() lane.Service {
	return lane.CreateService(
		f.CreateNotificationService(),
		f.RepositoryFactory.CreateRepository(),
		f.Logger)
}

// CreateCardService creates new service instance
func (f *ServiceFactory) CreateCardService() card.Service {
	return card.CreateService(
		f.CreateNotificationService(),
		f.RepositoryFactory.CreateRepository(),
		f.Logger)
}

// CreateCommandService creates new service instance
func (f *ServiceFactory) CreateCommandService() command.Service {
	return command.CreateService(
		f.CreateBoardService(),
		f.CreateLaneService(),
		f.CreateCardService(),
		f.Logger)
}

// CreateNotificationService creates new service instance
func (f *ServiceFactory) CreateNotificationService() notification.Service {
	return notification.CreateService(f.Publisher, f.Logger)
}
