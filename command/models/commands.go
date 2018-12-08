package command

import (
	"github.com/dmibod/kanban/kernel"
)

type Type int

const (
	UpdateCard Type = iota
	RemoveCard
	ExcludeCard
	InsertCard
)

type Command struct {
	Id      Id                `json:"id"`
	Type    Type              `json:"type"`
	Payload map[string]string `json:"payload"`
}
