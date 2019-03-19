package card

// Repository for card entity
type Repository interface {
	Create(*Entity) error
	Update(*Entity) error
	Delete(*Entity) error
}
