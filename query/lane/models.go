package lane

// Lane model
type Lane struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Layout      string `json:"layout"`
}
