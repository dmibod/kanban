package lane

// Repository for lane entity
type Repository interface {
	Create(*Entity) error
	Update(*Entity) error
	Delete(*Entity) error
}
