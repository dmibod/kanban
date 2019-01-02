package shared

import (
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db"
)

// CreateServiceFactory creates new instance
func CreateServiceFactory(f db.RepositoryFactory) *services.Factory {
	return services.CreateFactory(f, CreateLogger("[SERVICE] ", true))
}

// CreateRepositoryFactory creates repository factory
func CreateRepositoryFactory(s mongo.OperationExecutor) db.RepositoryFactory {
	return persistence.CreateFactory(s,	CreateLogger("[.MONGO.]", true))
}

// CreateDatabaseServices creates database services
func CreateDatabaseServices() (mongo.OperationExecutor, mongo.SessionProvider) {
	return persistence.CreateService(CreateLogger("[BRK.MGO] ", true))
}
