package query

// Board model
type Board struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Layout string `json:"layout,omitempty"`
	Owner  string `json:"owner,omitempty"`
	Shared bool   `json:"shared,omitempty"`
}

// Lane model
type Lane struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Type   string `json:"type"`
	Layout string `json:"layout"`
}

// Card model
type Card struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
