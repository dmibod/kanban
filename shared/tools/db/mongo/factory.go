package mongo

import (
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
)

var _ db.Factory = (*Factory)(nil)

// Factory declares repository factory
type Factory struct {
	executor OperationExecutor
	db       string
	logger   logger.Logger
}

// CreateFactory creates new repository factory
func CreateFactory(opts ...Option) *Factory {

	var options Options

	for _, o := range opts {
		o(&options)
	}

	l := options.logger

	if l == nil {
		l = &noop.Logger{}
	}

	return &Factory{
		logger:   l,
		db:       options.db,
		executor: options.executor,
	}
}

// CreateRepository creates new repository
func (f *Factory) CreateRepository(col string, instance db.InstanceFactory) db.Repository {
	return &Repository{
		executor: f.executor,
		instance: instance,
		ctx:      CreateOperationContext(f.db, col),
		logger:   f.logger,
	}
}
