package kernel

// CommandType type
type CommandType int

const (
	// UpdateCardCommand type 0
	UpdateCardCommand CommandType = CommandType(iota)
	// RemoveCardCommand type 1
	RemoveCardCommand
	// UpdateLaneCommand type 2
	UpdateLaneCommand
	// RemoveLaneCommand type 3
	RemoveLaneCommand
	// ExcludeChildCommand type 4
	ExcludeChildCommand
	// AppendChildCommand type 5
	AppendChildCommand
	// InsertBeforeCommand type 6
	InsertBeforeCommand
	// InsertAfterCommand type 7
	InsertAfterCommand
	// LayoutBoardCommand type 8
	LayoutBoardCommand
	// 9
	LayoutLaneCommand
	// 10
	DescribeBoardCommand
	// 11
	DescribeLaneCommand
	// 12
	DescribeCardCommand
	// 13
	UpdateBoardCommand
)

// Command type
type Command struct {
	ID      ID                `json:"id"`
	BoardID ID                `json:"board_id"`
	Type    CommandType       `json:"type"`
	Payload map[string]string `json:"payload"`
}
