package mongo

import (
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

var _ db.RepositoryFactory = (*repositoryFactory)(nil)

type repositoryFactory struct {
	logger.Logger
	OperationExecutor
	db string
}

// CreateRepositoryFactory creates repository factory
func CreateRepositoryFactory(db string, e OperationExecutor, l logger.Logger) db.RepositoryFactory {
	if l == nil {
		l = &noop.Logger{}
	}

	return &repositoryFactory{
		Logger:            l,
		db:                db,
		OperationExecutor: e,
	}
}

// CreateRepository creates new repository
func (f *repositoryFactory) CreateRepository(col string, r db.RepositoryEntity) db.Repository {
	return &repository{
		OperationExecutor: f.OperationExecutor,
		Logger:            f.Logger,
		RepositoryEntity:  r,
		db:                f.db,
		col:               col,
	}
}
