package kernel

type Id string

// CommandType type
type CommandType int

const (
	UpdateCard CommandType = CommandType(iota)
	RemoveCard
	ExcludeChild
	InsertBefore
	InsertAfter
	AppendChild
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
	RefreshCard NotificationType = NotificationType(iota)
	RefreshLane
	RefreshBoard
)

// Notification type
type Notification struct {
	Context Id               `json:"context"`
	ID      Id               `json:"id"`
	Type    NotificationType `json:"type"`
}
