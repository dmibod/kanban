package card

import (
	"github.com/dmibod/kanban/shared/kernel"
)

// CreatedEvent type
type CreatedEvent struct {
	Entity
}

// DeletedEvent type
type DeletedEvent struct {
	Entity
}

// NameChangedEvent type
type NameChangedEvent struct {
	ID       kernel.ID
	OldValue string
	NewValue string
}

// DescriptionChangedEvent type
type DescriptionChangedEvent struct {
	ID       kernel.ID
	OldValue string
	NewValue string
}
