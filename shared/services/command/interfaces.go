package command

import (
	"context"
	"errors"
	"github.com/dmibod/kanban/shared/kernel"
)

// Errors
var (
	ErrInvalidCommandType = errors.New("svc: invalid command type")
	ErrInvalidPayload     = errors.New("svc: invalid payload")
)

// Service interface
type Service interface {
	Execute(context.Context, kernel.Command) error
}
