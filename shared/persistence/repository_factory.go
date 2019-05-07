package persistence

import (
	"github.com/dmibod/kanban/shared/persistence/board"
	"github.com/dmibod/kanban/shared/persistence/card"
	"github.com/dmibod/kanban/shared/persistence/lane"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// RepositoryFactory type
type RepositoryFactory struct {
	factory *mongo.RepositoryFactory
}

// CreateBoardRepository method
func (f RepositoryFactory) CreateBoardRepository() board.Repository {
	return board.CreateRepository(f.factory.CreateRepository("boards"))
}

// CreateLaneRepository method
func (f RepositoryFactory) CreateLaneRepository() lane.Repository {
	return lane.CreateRepository(f.factory.CreateRepository("boards"))
}

// CreateCardRepository method
func (f RepositoryFactory) CreateCardRepository() card.Repository {
	return card.CreateRepository(f.factory.CreateRepository("boards"))
}

// CreateRepositoryFactory instance
func CreateRepositoryFactory(executor mongo.OperationExecutor, logger logger.Logger) RepositoryFactory {
	return RepositoryFactory{factory: mongo.CreateRepositoryFactory("kanban", executor, logger)}
}
