package mongo

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

// DatabaseCommand declares database command
type DatabaseCommand struct {
	db  string
	col string
}

// CreateDatabaseCommand creates DatabaseCommand
func CreateDatabaseCommand(db string, col string) *DatabaseCommand {
	return &DatabaseCommand{
		db:  db,
		col: col,
	}
}

// DatabaseCommandHandler declares DatabaseCommand handler
type DatabaseCommandHandler func(*mongo.Collection) error

// DatabaseCommandExecutor executes command
type DatabaseCommandExecutor interface {
	 Exec(*DatabaseCommand, DatabaseCommandHandler) error
}