package shared

import (
	"context"
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
	return persistence.CreateRepositoryFactory(e, CreateLogger("[.MONGO.]", true))
}

// CreateOperationExecutor instance
func CreateOperationExecutor(p mongo.SessionProvider) mongo.OperationExecutor {
	return persistence.CreateOperationExecutor(p, CreateLogger("[BRK.MGO] ", true))
}

// CreateContextFactory instance
func CreateContextFactory(p mongo.SessionProvider) mongo.ContextFactory {
	return mongo.CreateContextFactory(p, CreateLogger("[CTXFACT] ", true))
}

// CreateSessionProvider instance
func CreateSessionProvider(f mongo.SessionFactory) mongo.SessionProvider {
	return mongo.CreateSessionProvider(f, CreateLogger("[PROVIDR] ", true))
}

// CreateCopySessionProvider instance
func CreateCopySessionProvider(p mongo.SessionProvider) mongo.SessionProvider {
	return mongo.CreateCopySessionProvider(p, CreateLogger("[CPYPROV] ", true))
}

// CreateContextSessionProvider instance
func CreateContextSessionProvider(c context.Context) mongo.SessionProvider {
	return mongo.CreateContextSessionProvider(c, CreateLogger("[CTXPROV] ", true))
}

// CreateSessionFactory instance
func CreateSessionFactory() mongo.SessionFactory {
	return persistence.CreateSessionFactory(mongo.CreateSessionFactory(mongo.WithLogger(CreateLogger("[SESSION] ", true))), CreateLogger("[BRK.SES] ", true))
}
