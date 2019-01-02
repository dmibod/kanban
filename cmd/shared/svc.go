package shared

import (
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

// CreateServiceFactory instance
func CreateServiceFactory(f db.RepositoryFactory) *services.ServiceFactory {
	return services.CreateServiceFactory(f, CreateLogger("[SERVICE] ", true))
}

// CreateRepositoryFactory instance
func CreateRepositoryFactory(e mongo.OperationExecutor) db.RepositoryFactory {
	return persistence.CreateFactory(e, CreateLogger("[.MONGO.]", true))
}

// CreateExecutor instance
func CreateExecutor(f mongo.SessionFactory) mongo.OperationExecutor {
	return persistence.CreateExecutor(f, CreateLogger("[BRK.MGO] ", true))
}
