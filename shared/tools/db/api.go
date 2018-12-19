package db

// InstanceFactory declares entity instance factory
type InstanceFactory func() interface{}

// Factory declares repository factory
type Factory interface {
	CreateRepository(string, InstanceFactory) Repository
}

// Visitor desclares entity visitor
type Visitor func(interface{})

// Repository declares entity repository
type Repository interface {
	Create(interface{}) (string, error)

	FindByID(string) (interface{}, error)

	Find(interface{}, Visitor) error

	Count(interface{}) (int, error)

	Update(interface{}) error

	Remove(string) error
}
