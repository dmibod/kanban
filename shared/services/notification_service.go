package services

import (
	"encoding/json"

	"github.com/dmibod/kanban/shared/domain"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// NotificationService interface
type NotificationService interface {
	Execute(func(domain.EventRegistry) error) error
}

type notificationService struct {
	logger.Logger
	message.Publisher
}

// CreateNotificationService instance
func CreateNotificationService(p message.Publisher, l logger.Logger) NotificationService {
	return &notificationService{
		Publisher: p,
		Logger:    l,
	}
}

func (s *notificationService) Execute(handler func(domain.EventRegistry) error) error {
	if handler == nil {
		return nil
	}

	eventManager := domain.CreateEventManager()

	err := handler(eventManager)
	if err != nil {
		s.Errorln(err)
		return err
	}

	listener := &eventHandler{
		notifications: []kernel.Notification{},
		Logger:        s.Logger,
	}

	eventManager.Listen(listener)
	eventManager.Fire()

	err = listener.publish(s.Publisher)
	if err != nil {
		s.Errorln(err)
		return err
	}

	return nil
}

type eventHandler struct {
	logger.Logger
	notifications []kernel.Notification
}

func (n *eventHandler) Handle(event interface{}) {
	if event == nil {
		return
	}

	n.Debugf("domain event: %+v\n", event)

	if n.handleBoardEvent(event) {
		return
	}

	if n.handleLaneEvent(event) {
		return
	}

	n.handleCardEvent(event)
}

func (n *eventHandler) handleBoardEvent(event interface{}) bool {
	var notification kernel.Notification

	switch e := event.(type) {
	case domain.BoardCreatedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case domain.BoardDeletedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RemoveBoardNotification,
		}
	case domain.BoardNameChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case domain.BoardDescriptionChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case domain.BoardLayoutChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case domain.BoardSharedChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case domain.BoardChildAppendedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ChildID,
			Type:    kernel.RefreshBoardNotification,
		}
	case domain.BoardChildRemovedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ChildID,
			Type:    kernel.RefreshBoardNotification,
		}
	default:
		return false
	}

	for _, i := range n.notifications {
		if i.IsEqual(notification) {
			return true
		}
	}

	n.notifications = append(n.notifications, notification)

	return true
}

func (n *eventHandler) handleLaneEvent(event interface{}) bool {
	var notification kernel.Notification

	switch e := event.(type) {
	case domain.LaneCreatedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case domain.LaneDeletedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RemoveLaneNotification,
		}
	case domain.LaneNameChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case domain.LaneDescriptionChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case domain.LaneLayoutChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case domain.LaneChildAppendedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ChildID,
			Type:    kernel.RefreshLaneNotification,
		}
	case domain.LaneChildRemovedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ChildID,
			Type:    kernel.RefreshLaneNotification,
		}
	default:
		return false
	}

	for _, i := range n.notifications {
		if i.IsEqual(notification) {
			return true
		}
	}

	n.notifications = append(n.notifications, notification)

	return true
}

func (n *eventHandler) handleCardEvent(event interface{}) bool {
	var notification kernel.Notification

	switch e := event.(type) {
	case domain.CardCreatedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshCardNotification,
		}
	case domain.CardDeletedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RemoveCardNotification,
		}
	case domain.CardNameChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshCardNotification,
		}
	case domain.CardDescriptionChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshCardNotification,
		}
	default:
		return false
	}

	for _, i := range n.notifications {
		if i.IsEqual(notification) {
			return true
		}
	}

	n.notifications = append(n.notifications, notification)

	return true
}

func (n *eventHandler) publish(publisher message.Publisher) error {
	if len(n.notifications) == 0 {
		return nil
	}

	message, err := json.Marshal(n.notifications)
	if err != nil {
		return err
	}

	n.Debugf("publish notifications: %+v\n", n.notifications)

	return publisher.Publish(message)
}
