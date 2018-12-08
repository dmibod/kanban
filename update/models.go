package update

import (
	"github.com/dmibod/kanban/kernel"
)

type Card struct {
	Id kernel.Id `json:"id";omitempty;bson:"_id,omitempty"`
	Name string `json:"name";omitempty;bson:"name"`
}
