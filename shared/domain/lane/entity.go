package lane

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// Entity type
type Entity struct {
	ID          kernel.ID
	Kind        string
	Name        string
	Description string
	Layout      string
	Children    []kernel.ID
}
