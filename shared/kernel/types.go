package kernel

// ID type
type ID string

// EmptyID value
var EmptyID = ID("")

// IsValid id
func (id ID) IsValid() bool {
	return id != EmptyID
}

// String converts ID to string
func (id ID) String() string {
	return string(id)
}

// WithSet builds MemberID
func (id ID) WithSet(setID ID) MemberID {
	return MemberID{ID: id, SetID: setID}
}

// WithID builds MemberID
func (setID ID) WithID(id ID) MemberID {
	return MemberID{ID: id, SetID: setID}
}

// MemberID represents complex ID
type MemberID struct {
	SetID ID
	ID    ID
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
	BoardID ID                `json:"board_id"`
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
	ID      ID               `json:"id"`
	BoardID ID               `json:"board_id"`
	Type    NotificationType `json:"type"`
}

// IsEqual notification
func (n Notification) IsEqual(notification Notification) bool {
	return n.Type == notification.Type && n.BoardID == notification.BoardID && n.ID == notification.ID
}

// Layout
const (
	HLayout = "H"
	VLayout = "V"
)

// Lane Kind
const (
	LKind = "L"
	CKind = "C"
)
