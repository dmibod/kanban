package persistence

import (
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// CreateFactory creates new factory
func CreateFactory(e mongo.OperationExecutor, l logger.Logger) db.RepositoryFactory {
	return mongo.CreateFactory("kanban", e, l)
}
