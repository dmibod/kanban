package kernel

// CommandType type
type CommandType int

const (
	// UpdateCardCommand type
	UpdateCardCommand CommandType = CommandType(iota)
	// RemoveCardCommand type
	RemoveCardCommand
	// ExcludeChildCommand type
	ExcludeChildCommand
	// AppendChildCommand type
	AppendChildCommand
	// InsertBeforeCommand type
	InsertBeforeCommand
	// InsertAfterCommand type
	InsertAfterCommand
	// LayoutBoardCommand type
	LayoutBoardCommand
)

// Command type
type Command struct {
	ID      ID                `json:"id"`
	BoardID ID                `json:"board_id"`
	Type    CommandType       `json:"type"`
	Payload map[string]string `json:"payload"`
}
