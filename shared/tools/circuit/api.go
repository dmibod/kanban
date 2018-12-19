package circuit

type Subject interface {
	Execute() (interface{}, error)
}