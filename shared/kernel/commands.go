package kernel

// CommandType type
type CommandType int

const (
	// UpdateCardCommand type
	UpdateCardCommand CommandType = CommandType(iota)
	// RemoveCardCommand type
	RemoveCardCommand
	// UpdateLaneCommand type
	UpdateLaneCommand
	// RemoveLaneCommand type
	RemoveLaneCommand
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

	LayoutLaneCommand

	DescribeBoardCommand

	DescribeLaneCommand

	DescribeCardCommand
)

// Command type
type Command struct {
	ID      ID                `json:"id"`
	BoardID ID                `json:"board_id"`
	Type    CommandType       `json:"type"`
	Payload map[string]string `json:"payload"`
}
