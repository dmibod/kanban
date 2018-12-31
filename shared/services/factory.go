package services

import (
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// Factory creates service instances
type Factory struct {
	logger.Logger
	db.RepositoryFactory
}

// CreateFactory creates service factory
func CreateFactory(l logger.Logger, f db.RepositoryFactory) *Factory {
	return &Factory{
		Logger:            l,
		RepositoryFactory: f,
	}
}

// CreateCardService creates new service instance
func (f *Factory) CreateCardService() CardService {
	return &cardService{
		Logger:     f.Logger,
		Repository: persistence.CreateCardRepository(f.RepositoryFactory),
	}
}

// CreateBoardService creates new service instance
func (f *Factory) CreateBoardService() BoardService {
	return &boardService{
		Logger:     f.Logger,
		Repository: persistence.CreateBoardRepository(f.RepositoryFactory),
	}
}
