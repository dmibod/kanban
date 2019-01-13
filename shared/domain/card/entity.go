package card

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// Entity type
type Entity struct {
	ID          kernel.ID
	Name        string
	Description string
}
