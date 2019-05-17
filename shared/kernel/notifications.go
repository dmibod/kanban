package kernel

// NotificationType type
type NotificationType int

const (
	// RefreshCardNotification type
	RefreshCardNotification NotificationType = NotificationType(iota)
	// RefreshLaneNotification type
	RefreshLaneNotification
	// RefreshBoardNotification type
	RefreshBoardNotification
	// RemoveCardNotification type
	RemoveCardNotification
	// RemoveLaneNotification type
	RemoveLaneNotification
	// RemoveBoardNotification type
	RemoveBoardNotification
	// CreateCardNotification type
	CreateCardNotification
	// CreateLaneNotification type
	CreateLaneNotification
	// CreateBoardNotification type
	CreateBoardNotification
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
