package shared

import (
	"context"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"os"
)

const mongoURLEnvVar = "MGO_URL"
const mongoUSREnvVar = "MONGO_USER"
const mongoPWDEnvVar = "MONGO_PWD"

// CreateServiceFactory instance
func CreateServiceFactory(f persistence.RepositoryFactory) *services.ServiceFactory {
	return services.CreateServiceFactory(f, message.CreatePublisher("notification"), CreateLogger("[SERVICE] "))
}

// CreateRepositoryFactory instance
func CreateRepositoryFactory(e mongo.OperationExecutor) persistence.RepositoryFactory {
	return persistence.CreateRepositoryFactory(e, CreateLogger("[REPOFAC]"))
}

// CreateOperationExecutor instance
func CreateOperationExecutor(p mongo.SessionProvider) mongo.OperationExecutor {
	return persistence.CreateOperationExecutor(p, CreateLogger("[OPREXEC] "))
	//return mongo.CreateExecutor(p, CreateLogger("[OPREXEC] "))
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

func getMongoURLOrDefault(defURL string) string {
	url := os.Getenv(mongoURLEnvVar)

	if url == "" {
		return defURL
	}

	return url
}

func getMongoUserOrDefault(defUser string) string {
	user := os.Getenv(mongoUSREnvVar)

	if user == "" {
		return defUser
	}

	return user
}

func getMongoPasswordOrDefault(defPassword string) string {
	password := os.Getenv(mongoPWDEnvVar)

	if password == "" {
		return defPassword
	}

	return password
}

// CreateSessionFactory instance
func CreateSessionFactory() mongo.SessionFactory {
	sf := mongo.CreateSessionFactory(
		mongo.WithLogger(CreateLogger("[SESSFAC] ")),
		mongo.WithURL(getMongoURLOrDefault("")),
		mongo.WithUser(getMongoUserOrDefault("")),
		mongo.WithPassword(getMongoPasswordOrDefault("")))

	return persistence.CreateSessionFactory(sf, CreateLogger("[BRK.SES] "))
}
