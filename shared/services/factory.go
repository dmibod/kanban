package services

import (
	"context"

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

// CreateCardService creates new CardService instance
func (f *Factory) CreateCardService(ctx context.Context) CardService {
	return &cardService{
		Context:           ctx,
		Logger:            f.Logger,
		RepositoryFactory: f.RepositoryFactory,
	}
}
