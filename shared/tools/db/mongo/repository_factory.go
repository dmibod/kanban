package mongo

import (
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

var _ db.RepositoryFactory = (*repositoryFactory)(nil)

type repositoryFactory struct {
	executor OperationExecutor
	db       string
	logger   logger.Logger
}

// CreateFactory creates repository factory
func CreateFactory(db string, e OperationExecutor, l logger.Logger) db.RepositoryFactory {
	if l == nil {
		l = &noop.Logger{}
	}

	return &repositoryFactory{
		logger:   l,
		db:       db,
		executor: e,
	}
}

// CreateRepository creates new repository
func (f *repositoryFactory) CreateRepository(col string, instanceFactory db.InstanceFactory, instanceIdentity db.InstanceIdentity) db.Repository {
	return &repository{
		executor:         f.executor,
		instanceFactory:  instanceFactory,
		instanceIdentity: instanceIdentity,
		db:               f.db,
		col:              col,
		logger:           f.logger,
	}
}
