package services

import (
	"encoding/json"

	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/domain/card"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/domain/lane"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// NotificationService interface
type NotificationService interface {
	Execute(func(event.Registry) error) error
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

func (s *notificationService) Execute(handler func(event.Registry) error) error {
	if handler == nil {
		return nil
	}

	manager := event.CreateEventManager()

	err := handler(manager)
	if err != nil {
		s.Errorln(err)
		return err
	}

	listener := &eventHandler{
		notifications: []kernel.Notification{},
		Logger:        s.Logger,
	}

	manager.Listen(listener)
	manager.Fire()

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
	case board.CreatedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.DeletedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RemoveBoardNotification,
		}
	case board.NameChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.DescriptionChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.LayoutChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.SharedChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.ChildAppendedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ChildID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.ChildRemovedEvent:
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
	case lane.CreatedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case lane.DeletedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RemoveLaneNotification,
		}
	case lane.NameChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case lane.DescriptionChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case lane.LayoutChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case lane.ChildAppendedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ChildID,
			Type:    kernel.RefreshLaneNotification,
		}
	case lane.ChildRemovedEvent:
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
	case card.CreatedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshCardNotification,
		}
	case card.DeletedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RemoveCardNotification,
		}
	case card.NameChangedEvent:
		notification = kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshCardNotification,
		}
	case card.DescriptionChangedEvent:
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
