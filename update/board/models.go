package board

// Board api
type Board struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Layout string `json:"layout,omitempty"`
	Owner  string `json:"owner,omitempty"`
	Shared bool   `json:"shared,omitempty"`
}
