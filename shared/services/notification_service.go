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
	Listen(event.Bus)
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

// Listen doamin events
func (s *notificationService) Listen(bus event.Bus) {
	if bus != nil {
		bus.Listen(s)
	}
}

// Handle doamin event
func (s *notificationService) Handle(event interface{}) {
	if event == nil {
		return
	}

	s.Debugf("domain event: %+v\n", event)

	if !s.handleBoardEvent(event) {
		if !s.handleLaneEvent(event) {
			s.handleCardEvent(event)
		}
	}
}

func (s *notificationService) handleBoardEvent(event interface{}) bool {
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

	s.publish(notification)

	return true
}

func (s *notificationService) handleLaneEvent(event interface{}) bool {
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

	s.publish(notification)

	return true
}

func (s *notificationService) handleCardEvent(event interface{}) bool {
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

	s.publish(notification)

	return true
}

func (s *notificationService) publish(notification kernel.Notification) {
	message, err := json.Marshal([]kernel.Notification{notification})
	if err != nil {
		s.Errorln(err)
		return
	}

	s.Debugf("publish notification: %+v\n", &notification)

	err = s.Publisher.Publish(message)
	if err != nil {
		s.Errorln(err)
		return
	}
}
