package services

import (
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
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
func (f *ServiceFactory) CreateBoardService() BoardService {
	return CreateBoardService(
		f.CreateNotificationService(),
		persistence.CreateBoardRepository(f.RepositoryFactory),
		f.Logger)
}

// CreateLaneService creates new service instance
func (f *ServiceFactory) CreateLaneService() LaneService {
	return CreateLaneService(
		f.CreateNotificationService(),
		persistence.CreateLaneRepository(f.RepositoryFactory),
		persistence.CreateBoardRepository(f.RepositoryFactory),
		f.Logger)
}

// CreateCardService creates new service instance
func (f *ServiceFactory) CreateCardService() CardService {
	return CreateCardService(
		f.CreateNotificationService(),
		persistence.CreateCardRepository(f.RepositoryFactory),
		persistence.CreateLaneRepository(f.RepositoryFactory),
		f.Logger)
}

// CreateCommandService creates new service instance
func (f *ServiceFactory) CreateCommandService() CommandService {
	return &commandService{
		Logger:       f.Logger,
		boardService: f.CreateBoardService(),
		laneService:  f.CreateLaneService(),
		cardService:  f.CreateCardService(),
	}
}

// CreateNotificationService creates new service instance
func (f *ServiceFactory) CreateNotificationService() NotificationService {
	return CreateNotificationService(f.Publisher, f.Logger)
}
