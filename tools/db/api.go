package db

type Visitor func(interface{})

type Repository interface {
	
	Create(interface{}) (string, error)
	
	FindById(string) (interface{}, error)
	
	Find(interface{}, Visitor) error

	Count(interface{}) (int, error)
	
	Update(interface{}) error
	
	Remove(string) error
}