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

	n.handleBoardEvent(event)
}

func (n *eventHandler) handleBoardEvent(event interface{}) bool {
	var notification *kernel.Notification

	switch e := event.(type) {
	case *domain.BoardNameChangedEvent:
		notification = &kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case *domain.BoardDescriptionChangedEvent:
		notification = &kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case *domain.BoardLayoutChangedEvent:
		notification = &kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case *domain.BoardSharedChangedEvent:
		notification = &kernel.Notification{
			Context: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case *domain.BoardChildAppendedEvent:
		notification = &kernel.Notification{
			Context: e.ID,
			ID:      e.ChildID,
			Type:    kernel.RefreshBoardNotification,
		}
	case *domain.BoardChildRemovedEvent:
		notification = &kernel.Notification{
			Context: e.ID,
			ID:      e.ChildID,
			Type:    kernel.RefreshBoardNotification,
		}
	default:
		return false
	}

	if notification != nil {
		for _, i := range n.notifications {
			if i.Type == notification.Type && i.Context == notification.Context && i.ID == notification.ID {
				return true
			}
		}

		n.notifications = append(n.notifications, *notification)
	}

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
