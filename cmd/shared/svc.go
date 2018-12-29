package shared

import (
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
)

// CreateServiceFactory creates new instance
func CreateServiceFactory() *services.Factory {
	return services.CreateFactory(createLogger("[SERVICE] ", true), createRepositoryFactory())
}

func createRepositoryFactory() db.RepositoryFactory {
	return persistence.CreateFactory(
		persistence.CreateService(createLogger("[BRK.MGO]", true)),
		createLogger("[MONGO..]", true))
}

func createLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
