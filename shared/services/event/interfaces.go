package event

import (
	"context"
	"github.com/dmibod/kanban/shared/domain/event"
)

// Service interface
type Service interface {
	Execute(context.Context, func(event.Bus) error) error
}
