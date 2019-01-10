package db

import (
	"context"
)

// EntityFactory interface
type EntityFactory interface {
	CreateInstance() interface{}
}

// EntityIdentity interface
type EntityIdentity interface {
	GetID(interface{}) string
}

// RepositoryEntity interface
type RepositoryEntity interface {
	EntityFactory
	EntityIdentity
}

// RepositoryFactory interface
type RepositoryFactory interface {
	CreateRepository(string, RepositoryEntity) Repository
}

// EntityVisitor entity visitor, must return true to stop iteration
type EntityVisitor func(interface{}) error

// Repository interface
type Repository interface {
	Create(context.Context, interface{}) (string, error)

	FindByID(context.Context, string) (interface{}, error)

	Find(context.Context, interface{}, EntityVisitor) error

	Count(context.Context, interface{}) (int, error)

	Update(context.Context, interface{}) error

	Remove(context.Context, string) error
}
