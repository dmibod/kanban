package domain

import "github.com/dmibod/kanban/shared/kernel"

// Saver interface
type Saver interface {
	Save() error
}

// Repository interface
type Repository interface {
	Fetch(kernel.Id) (interface{}, error)
	Persist(interface{}) (kernel.Id, error)
}
