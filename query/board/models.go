package board

// ListModel type
type ListModel struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
	State       string `json:"state,omitempty"`
	Shared      bool   `json:"shared,omitempty"`
}

// Model type
type Model struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Layout      string      `json:"layout,omitempty"`
	Owner       string      `json:"owner,omitempty"`
	State       string      `json:"state,omitempty"`
	Shared      bool        `json:"shared,omitempty"`
	Lanes       []LaneModel `json:"lanes,omitempty"`
}

// LaneModel type
type LaneModel struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Type        string      `json:"type"`
	Layout      string      `json:"layout"`
	Lanes       []LaneModel `json:"lanes,omitempty"`
	Cards       []CardModel `json:"cards,omitempty"`
}

// CardModel type
type CardModel struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
