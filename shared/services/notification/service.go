package notification

import (
	"context"
	"encoding/json"

	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/domain/card"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/domain/lane"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type service struct {
	logger.Logger
	message.Publisher
}

// CreateService instance
func CreateService(p message.Publisher, l logger.Logger) Service {
	return &service{
		Publisher: p,
		Logger:    l,
	}
}

// Listen doamin events
func (s *service) Listen(bus event.Bus) {
	if bus != nil {
		bus.Listen(s)
	}
}

// Handle doamin event
func (s *service) Handle(ctx context.Context, event interface{}) error {
	if event == nil {
		return nil
	}

	s.Debugf("domain event: %+v\n", event)

	if !s.handleBoardEvent(event) {
		if !s.handleLaneEvent(event) {
			s.handleCardEvent(event)
		}
	}

	return nil
}

func (s *service) handleBoardEvent(event interface{}) bool {
	var notification kernel.Notification

	switch e := event.(type) {
	case board.CreatedEvent:
		notification = kernel.Notification{
			BoardID: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.DeletedEvent:
		notification = kernel.Notification{
			BoardID: e.ID,
			ID:      e.ID,
			Type:    kernel.RemoveBoardNotification,
		}
	case board.NameChangedEvent:
		notification = kernel.Notification{
			BoardID: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.DescriptionChangedEvent:
		notification = kernel.Notification{
			BoardID: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.LayoutChangedEvent:
		notification = kernel.Notification{
			BoardID: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.SharedChangedEvent:
		notification = kernel.Notification{
			BoardID: e.ID,
			ID:      e.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.ChildAppendedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	case board.ChildRemovedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RefreshBoardNotification,
		}
	default:
		return false
	}

	s.publish(notification)

	return true
}

func (s *service) handleLaneEvent(event interface{}) bool {
	var notification kernel.Notification

	switch e := event.(type) {
	case lane.CreatedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case lane.DeletedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RemoveLaneNotification,
		}
	case lane.NameChangedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case lane.DescriptionChangedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case lane.LayoutChangedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case lane.ChildAppendedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	case lane.ChildRemovedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RefreshLaneNotification,
		}
	default:
		return false
	}

	s.publish(notification)

	return true
}

func (s *service) handleCardEvent(event interface{}) bool {
	var notification kernel.Notification

	switch e := event.(type) {
	case card.CreatedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RefreshCardNotification,
		}
	case card.DeletedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RemoveCardNotification,
		}
	case card.NameChangedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RefreshCardNotification,
		}
	case card.DescriptionChangedEvent:
		notification = kernel.Notification{
			BoardID: e.ID.SetID,
			ID:      e.ID.ID,
			Type:    kernel.RefreshCardNotification,
		}
	default:
		return false
	}

	s.publish(notification)

	return true
}

func (s *service) publish(notification kernel.Notification) {
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
