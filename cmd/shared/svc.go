package shared

import (
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db"
)

// CreateServiceFactory creates new instance
func CreateServiceFactory() *services.Factory {
	return services.CreateFactory(
		CreateLogger("[SERVICE] ", true),
		createRepositoryFactory())
}

func createRepositoryFactory() db.RepositoryFactory {
	return persistence.CreateFactory(
		persistence.CreateService(CreateLogger("[BRK.MGO]", true)),
		CreateLogger("[.MONGO.]", true))
}
