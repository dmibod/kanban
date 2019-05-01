package card

// Aggregate interface
type Aggregate interface {
	Root() Entity
	Name(string) error
	Description(string) error
}
