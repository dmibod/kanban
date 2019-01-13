package board

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// Entity type
type Entity struct {
	ID          kernel.ID
	Owner       string
	Name        string
	Description string
	Layout      string
	Shared      bool
	Children    []kernel.ID
}
