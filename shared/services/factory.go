package services

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// Factory creates service instances
type Factory struct {
	logger  logger.Logger
	factory db.Factory
}

// CreateFactory creates service factory
func CreateFactory(l logger.Logger, f db.Factory) *Factory {
	return &Factory{
		logger:  l,
		factory: f,
	}
}

// CreateCardService creates new CardService instance
func (f *Factory) CreateCardService(c context.Context) CardService {
	return &cardService{
		ctx:     c,
		logger:  f.logger,
		factory: f.factory,
	}
}
