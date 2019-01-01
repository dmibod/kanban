package db

import (
	"context"
)

// InstanceFactory creates entity instance
type InstanceFactory func() interface{}

// InstanceIdentity gets an id from entity
type InstanceIdentity func(interface{}) string

// RepositoryFactory interface
type RepositoryFactory interface {
	CreateRepository(string, InstanceFactory, InstanceIdentity) Repository
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
