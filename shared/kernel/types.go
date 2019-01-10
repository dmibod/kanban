package kernel

// EmptyID value
var EmptyID = Id("")

// Id type
type Id string

// IsValid id
func (id Id) IsValid() bool {
	sid := string(id)
	return sid != ""
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
	ID      Id                `json:"id"`
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
	Context Id               `json:"context"`
	ID      Id               `json:"id"`
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
