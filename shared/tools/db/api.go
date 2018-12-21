package db

import (
	"context"
)

// InstanceFactory declares entity instance factory
type InstanceFactory func() interface{}

// RepositoryFactory interface
type RepositoryFactory interface {
	CreateRepository(context.Context, string, InstanceFactory) Repository
}

// EntityVisitor entity visitor, must return true to stop iteration
type EntityVisitor func(interface{}) bool

// Repository interface
type Repository interface {
	Create(interface{}) (string, error)

	FindByID(string) (interface{}, error)

	Find(interface{}, EntityVisitor) error

	Count(interface{}) (int, error)

	Update(interface{}) error

	Remove(string) error
}
