package error

import (
	"errors"
)

// Errors
var (
	ErrInvalidArgument = errors.New("agg: invalid argument")
	ErrInvalidID       = errors.New("agg: invalid id")
	ErrInvalidType     = errors.New("agg: invalid type")
)
