package persistence

import (
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// RepositoryFactory type
type RepositoryFactory struct {
	factory *mongo.RepositoryFactory
}

// CreateRepository func
func (f RepositoryFactory) CreateRepository() Repository {
	return Repository{repository: f.factory.CreateRepository("boards")}
}

// CreateRepositoryFactory instance
func CreateRepositoryFactory(executor mongo.OperationExecutor, logger logger.Logger) RepositoryFactory {
	return RepositoryFactory{factory: mongo.CreateRepositoryFactory("kanban", executor, logger)}
}
