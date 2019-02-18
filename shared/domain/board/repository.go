package board

// Repository for board entity
type Repository interface {
	Create(*Entity) error
	Update(*Entity) error
	Delete(*Entity) error
}
