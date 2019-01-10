package shared

import (
	"github.com/dmibod/kanban/shared/message"
	"context"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"os"
)

const mongoUrlEnvVar = "MGO_URL"

// CreateServiceFactory instance
func CreateServiceFactory(f db.RepositoryFactory) *services.ServiceFactory {
	return services.CreateServiceFactory(f, message.CreatePublisher("notification"), CreateLogger("[SERVICE] "))
}

// CreateRepositoryFactory instance
func CreateRepositoryFactory(e mongo.OperationExecutor) db.RepositoryFactory {
	return persistence.CreateRepositoryFactory(e, CreateLogger("[REPOFAC]"))
}

// CreateOperationExecutor instance
func CreateOperationExecutor(p mongo.SessionProvider) mongo.OperationExecutor {
	return persistence.CreateOperationExecutor(p, CreateLogger("[OPREXEC] "))
}

// CreateContextFactory instance
func CreateContextFactory(p mongo.SessionProvider) mongo.ContextFactory {
	return mongo.CreateContextFactory(p, CreateLogger("[CTXFACT] "))
}

// CreateSessionProvider instance
func CreateSessionProvider(f mongo.SessionFactory) mongo.SessionProvider {
	return mongo.CreateSessionProvider(f, CreateLogger("[SESPROV] "))
}

// CreateCopySessionProvider instance
func CreateCopySessionProvider(p mongo.SessionProvider) mongo.SessionProvider {
	return mongo.CreateCopySessionProvider(p, CreateLogger("[CPYPROV] "))
}

// CreateContextSessionProvider instance
func CreateContextSessionProvider(c context.Context) mongo.SessionProvider {
	return mongo.CreateContextSessionProvider(c, CreateLogger("[CTXPROV] "))
}

func getMongoUrlOrDefault(defUrl string) string {
	url := os.Getenv(mongoUrlEnvVar)

	if url == "" {
		return defUrl
	}

	return url
}

// CreateSessionFactory instance
func CreateSessionFactory() mongo.SessionFactory {
	sf := mongo.CreateSessionFactory(
		mongo.WithLogger(CreateLogger("[SESSFAC] ")),
		mongo.WithURL(getMongoUrlOrDefault("")))

	return persistence.CreateSessionFactory(sf, CreateLogger("[BRK.SES] "))
}
