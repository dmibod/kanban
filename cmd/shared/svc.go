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
func CreateExecutor(p mongo.SessionProvider) mongo.OperationExecutor {
	return persistence.CreateExecutor(p, CreateLogger("[BRK.MGO] ", true))
}

// CreateContextFactory instance
func CreateContextFactory(p mongo.SessionProvider) mongo.ContextFactory {
	return mongo.CreateContextFactory(p, CreateLogger("[CTXFACT] ", true))
}

// CreateSessionProvider instance
func CreateSessionProvider(f mongo.SessionFactory) mongo.SessionProvider {
	return mongo.CreateSessionProvider(f, CreateLogger("[PROVIDR] ", true))
}

// CreateSessionFactory instance
func CreateSessionFactory() mongo.SessionFactory {
	return persistence.CreateSessionFactory(mongo.CreateSessionFactory(mongo.WithLogger(CreateLogger("[SESSION] ", true))), CreateLogger("[BRK.SES] ", true))
}
