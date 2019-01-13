package persistence

import (
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// CreateRepositoryFactory instance
func CreateRepositoryFactory(e mongo.OperationExecutor, l logger.Logger) *mongo.RepositoryFactory {
	return mongo.CreateRepositoryFactory("kanban", e, l)
}
