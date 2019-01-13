package mongo

import (
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

// RepositoryFactory type
type RepositoryFactory struct {
	logger.Logger
	OperationExecutor
	db string
}

// CreateRepositoryFactory creates repository factory
func CreateRepositoryFactory(db string, e OperationExecutor, l logger.Logger) *RepositoryFactory {
	if l == nil {
		l = &noop.Logger{}
	}

	return &RepositoryFactory{
		Logger:            l,
		db:                db,
		OperationExecutor: e,
	}
}

// CreateRepository creates new repository
func (f *RepositoryFactory) CreateRepository(col string) *Repository {
	return &Repository{
		OperationExecutor: f.OperationExecutor,
		Logger:            f.Logger,
		db:                f.db,
		col:               col,
	}
}
