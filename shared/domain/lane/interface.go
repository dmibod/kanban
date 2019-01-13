package lane

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// Aggregate interface
type Aggregate interface {
	Root() Entity
	Name(string) error
	Description(string) error
	Layout(string) error
	AppendChild(kernel.ID) error
	RemoveChild(kernel.ID) error
}
