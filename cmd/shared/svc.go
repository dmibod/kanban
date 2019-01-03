package shared

import (
	"context"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

const debug = false

// CreateServiceFactory instance
func CreateServiceFactory(f db.RepositoryFactory) *services.ServiceFactory {
	return services.CreateServiceFactory(f, CreateLogger("[SERVICE] ", debug))
}

// CreateRepositoryFactory instance
func CreateRepositoryFactory(e mongo.OperationExecutor) db.RepositoryFactory {
	return persistence.CreateRepositoryFactory(e, CreateLogger("[REPOFAC]", debug))
}

// CreateOperationExecutor instance
func CreateOperationExecutor(p mongo.SessionProvider) mongo.OperationExecutor {
	return persistence.CreateOperationExecutor(p, CreateLogger("[OPREXEC] ", debug))
}

// CreateContextFactory instance
func CreateContextFactory(p mongo.SessionProvider) mongo.ContextFactory {
	return mongo.CreateContextFactory(p, CreateLogger("[CTXFACT] ", debug))
}

// CreateSessionProvider instance
func CreateSessionProvider(f mongo.SessionFactory) mongo.SessionProvider {
	return mongo.CreateSessionProvider(f, CreateLogger("[SESPROV] ", debug))
}

// CreateCopySessionProvider instance
func CreateCopySessionProvider(p mongo.SessionProvider) mongo.SessionProvider {
	return mongo.CreateCopySessionProvider(p, CreateLogger("[CPYPROV] ", debug))
}

// CreateContextSessionProvider instance
func CreateContextSessionProvider(c context.Context) mongo.SessionProvider {
	return mongo.CreateContextSessionProvider(c, CreateLogger("[CTXPROV] ", debug))
}

// CreateSessionFactory instance
func CreateSessionFactory() mongo.SessionFactory {
	return persistence.CreateSessionFactory(mongo.CreateSessionFactory(mongo.WithLogger(CreateLogger("[SESSFAC] ", debug))), CreateLogger("[BRK.SES] ", debug))
}
