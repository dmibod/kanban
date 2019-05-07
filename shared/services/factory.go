package services

import (
	"github.com/dmibod/kanban/shared/message"
	db "github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services/board"
	"github.com/dmibod/kanban/shared/services/card"
	"github.com/dmibod/kanban/shared/services/command"
	"github.com/dmibod/kanban/shared/services/event"
	"github.com/dmibod/kanban/shared/services/lane"
	"github.com/dmibod/kanban/shared/services/notification"
	"github.com/dmibod/kanban/shared/services/persistence"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// ServiceFactory creates service instances
type ServiceFactory struct {
	db.RepositoryFactory
	message.Publisher
	logger.Logger
}

// CreateServiceFactory creates service factory
func CreateServiceFactory(f db.RepositoryFactory, p message.Publisher, l logger.Logger) *ServiceFactory {
	return &ServiceFactory{
		RepositoryFactory: f,
		Publisher:         p,
		Logger:            l,
	}
}

// CreateBoardService creates new service instance
func (f *ServiceFactory) CreateBoardService() board.Service {
	return board.CreateService(
		f.CreateEventService(),
		f.RepositoryFactory.CreateBoardRepository(),
		f.Logger)
}

// CreateLaneService creates new service instance
func (f *ServiceFactory) CreateLaneService() lane.Service {
	return lane.CreateService(
		f.CreateEventService(),
		f.RepositoryFactory.CreateLaneRepository(),
		f.Logger)
}

// CreateCardService creates new service instance
func (f *ServiceFactory) CreateCardService() card.Service {
	return card.CreateService(
		f.CreateEventService(),
		f.RepositoryFactory.CreateCardRepository(),
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

// CreateEventService creates new service instance
func (f *ServiceFactory) CreateEventService() event.Service {
	return event.CreateService(f.CreateNotificationService(), f.CreatePersistenceService(), f.Logger)
}

// CreatePersistenceService creates new service instance
func (f *ServiceFactory) CreatePersistenceService() persistence.Service {
	return persistence.CreateService(
		f.RepositoryFactory.CreateBoardRepository(),
		f.RepositoryFactory.CreateLaneRepository(),
		f.RepositoryFactory.CreateCardRepository(),
		f.Logger)
}

// CreateNotificationService creates new service instance
func (f *ServiceFactory) CreateNotificationService() notification.Service {
	return notification.CreateService(f.Publisher, f.Logger)
}
