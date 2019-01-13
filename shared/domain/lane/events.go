package lane

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

// LayoutChangedEvent type
type LayoutChangedEvent struct {
	ID       kernel.ID
	OldValue string
	NewValue string
}

// ChildAppendedEvent type
type ChildAppendedEvent struct {
	ID      kernel.ID
	ChildID kernel.ID
}

// ChildRemovedEvent type
type ChildRemovedEvent struct {
	ID      kernel.ID
	ChildID kernel.ID
}
