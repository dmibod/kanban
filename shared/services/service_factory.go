package services

import (
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// ServiceFactory creates service instances
type ServiceFactory struct {
	db.RepositoryFactory
	logger.Logger
}

// CreateServiceFactory creates service factory
func CreateServiceFactory(f db.RepositoryFactory, l logger.Logger) *ServiceFactory {
	return &ServiceFactory{
		RepositoryFactory: f,
		Logger:            l,
	}
}

// CreateBoardService creates new service instance
func (f *ServiceFactory) CreateBoardService() BoardService {
	return &boardService{
		Logger:     f.Logger,
		Repository: persistence.CreateBoardRepository(f.RepositoryFactory),
	}
}

// CreateLaneService creates new service instance
func (f *ServiceFactory) CreateLaneService() LaneService {
	return &laneService{
		Logger:     f.Logger,
		Repository: persistence.CreateLaneRepository(f.RepositoryFactory),
	}
}

// CreateCardService creates new service instance
func (f *ServiceFactory) CreateCardService() CardService {
	return &cardService{
		Logger:         f.Logger,
		cardRepository: persistence.CreateCardRepository(f.RepositoryFactory),
		laneRepository: persistence.CreateLaneRepository(f.RepositoryFactory),
	}
}
