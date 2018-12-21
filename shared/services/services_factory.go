package services

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// Factory creates service instances
type Factory struct {
	logger            logger.Logger
	repositoryFactory db.RepositoryFactory
}

// CreateFactory creates service factory
func CreateFactory(l logger.Logger, f db.RepositoryFactory) *Factory {
	return &Factory{
		logger:            l,
		repositoryFactory: f,
	}
}

// CreateCardService creates new CardService instance
func (f *Factory) CreateCardService(ctx context.Context) CardService {
	return &cardService{
		ctx:               ctx,
		logger:            f.logger,
		repositoryFactory: f.repositoryFactory,
	}
}
