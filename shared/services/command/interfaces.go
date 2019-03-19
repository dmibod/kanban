package command

import (
	"github.com/dmibod/kanban/shared/kernel"
	"context"
	"errors"
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
