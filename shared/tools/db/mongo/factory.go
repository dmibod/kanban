package mongo

import (
	"github.com/dmibod/kanban/shared/tools/log/logger"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/log"
)

var _ db.Factory = (*Factory)(nil)

// Factory declares repository factory
type Factory struct {
	executor DatabaseCommandExecutor
	db       string
	logger   log.Logger
}

// CreateFactory creates new repository factory
func CreateFactory(opts ...Option) *Factory {

	var options Options

	for _, o := range opts {
		o(&options)
	}

	l := options.logger

	if l == nil {
		l = logger.New(logger.WithPrefix("[MONGO] "), logger.WithDebug(true))
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
		cmd:      CreateDatabaseCommand(f.db, col),
		logger:   f.logger,
	}
}
