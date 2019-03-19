package board

// ListModel type
type ListModel struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Shared      bool   `json:"shared,omitempty"`
}

// Model type
type Model struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Layout      string `json:"layout,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Shared      bool   `json:"shared,omitempty"`
}
