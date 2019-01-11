package kernel

// EmptyID value
var EmptyID = ID("")

// ID type
type ID string

// IsValid id
func (id ID) IsValid() bool {
	return id != EmptyID
}

// CommandType type
type CommandType int

const (
	UpdateCard CommandType = CommandType(iota)
	RemoveCard
	ExcludeChild
	AppendChild
	InsertBefore
	InsertAfter
	LayoutBoard
)

// Command type
type Command struct {
	ID      ID                `json:"id"`
	Type    CommandType       `json:"type"`
	Payload map[string]string `json:"payload"`
}

// NotificationType type
type NotificationType int

const (
	RefreshCardNotification NotificationType = NotificationType(iota)
	RefreshLaneNotification
	RefreshBoardNotification
	RemoveCardNotification
	RemoveLaneNotification
	RemoveBoardNotification
)

// Notification type
type Notification struct {
	Context ID               `json:"context"`
	ID      ID               `json:"id"`
	Type    NotificationType `json:"type"`
}

// Layout
const (
	HLayout = "H"
	VLayout = "V"
)

// Lane Type
const (
	LType = "L"
	CType = "C"
)
