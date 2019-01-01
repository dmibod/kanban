package services

import (
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// Factory creates service instances
type Factory struct {
	db.RepositoryFactory
	logger.Logger
}

// CreateFactory creates service factory
func CreateFactory(f db.RepositoryFactory, l logger.Logger) *Factory {
	return &Factory{
		RepositoryFactory: f,
		Logger:            l,
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
