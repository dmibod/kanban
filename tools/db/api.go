package db

type VisitFn func(interface{})

type Repository interface {
	
	Create(interface{}) (string, error)
	
	FindById(string) (interface{}, error)
	
	Find(interface{}, VisitFn) error

	Count(interface{}) (int, error)
	
	Update(interface{}) error
	
	Remove(string) error
}