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
	CreateRepository(context.Context, string, InstanceFactory, InstanceIdentity) Repository
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
