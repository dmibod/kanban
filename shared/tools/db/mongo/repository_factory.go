package mongo

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

var _ db.RepositoryFactory = (*repositoryFactory)(nil)

// RepositoryFactory repository factory
type repositoryFactory struct {
	executor OperationExecutor
	db       string
	logger   logger.Logger
}

// CreateFactory creates new repository factory
func CreateFactory(opts ...Option) db.RepositoryFactory {

	var options Options

	for _, o := range opts {
		o(&options)
	}

	l := options.logger

	if l == nil {
		l = &noop.Logger{}
	}

	return &repositoryFactory{
		logger:   l,
		db:       options.db,
		executor: options.executor,
	}
}

// CreateRepository creates new repository
func (f *repositoryFactory) CreateRepository(ctx context.Context, col string, instanceFactory db.InstanceFactory) db.Repository {
	return &Repository{
		executor:        f.executor,
		instanceFactory: instanceFactory,
		ctx:             CreateOperationContext(ctx, f.db, col),
		logger:          f.logger,
	}
}

// Options declares repository factory options
type Options struct {
	executor OperationExecutor
	logger   logger.Logger
	db       string
}

// Option is a closure which should initialize specific Options properties
type Option func(*Options)

// WithLogger initializes logger option
func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.logger = l
	}
}

// WithDatabase initializes db option
func WithDatabase(db string) Option {
	return func(o *Options) {
		o.db = db
	}
}

// WithExecutor initializes executor option
func WithExecutor(e OperationExecutor) Option {
	return func(o *Options) {
		o.executor = e
	}
}
