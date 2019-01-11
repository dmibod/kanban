package services

import (
	"context"
	"errors"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// Errors
var (
	ErrInvalidCommandType = errors.New("svc: invalid command type")
	ErrInvalidPayload     = errors.New("svc: invalid payload")
)

// CommandService interface
type CommandService interface {
	Execute(context.Context, kernel.Command) error
}

type commandService struct {
	logger.Logger
	boardService BoardService
	laneService  LaneService
	cardService  CardService
}

func (s *commandService) Execute(ctx context.Context, command kernel.Command) error {
	result := ErrInvalidCommandType

	switch command.Type {
	case kernel.InsertBefore: //todo
	case kernel.InsertAfter: //todo
	case kernel.AppendChild:
		if parentID, err := s.getID("parent_id", command.Payload); err != nil {
			result = err
		} else {
			return s.appendChild(ctx, command.ID, parentID)
		}
	case kernel.ExcludeChild:
		if parentID, err := s.getID("parent_id", command.Payload); err != nil {
			result = err
		} else {
			return s.excludeChild(ctx, command.ID, parentID)
		}
	case kernel.UpdateCard:
		if name, err := s.getString("name", command.Payload); err != nil {
			result = err
		} else {
			return s.updateCard(ctx, command.ID, name)
		}
	case kernel.RemoveCard:
		if parentID, err := s.getID("parent_id", command.Payload); err != nil {
			result = err
		} else {
			return s.removeCard(ctx, command.ID, parentID)
		}
	case kernel.LayoutBoard:
		if layout, err := s.getString("layout", command.Payload); err != nil {
			result = err
		} else {
			return s.layoutBoard(ctx, command.ID, layout)
		}
	}

	return result
}

func (s *commandService) insertBefore(ctx context.Context, id kernel.ID, relativeId kernel.ID) error {
	return nil
}

func (s *commandService) insertAfter(ctx context.Context, id kernel.ID, relativeId kernel.ID) error {
	return nil
}

func (s *commandService) appendChild(ctx context.Context, id kernel.ID, parentId kernel.ID) error {
	if _, err := s.boardService.GetByID(ctx, parentId); err == nil {
		return s.boardService.AppendChild(ctx, parentId, id)
	}
	return s.laneService.AppendChild(ctx, parentId, id)
}

func (s *commandService) excludeChild(ctx context.Context, id kernel.ID, parentId kernel.ID) error {
	if _, err := s.boardService.GetByID(ctx, parentId); err == nil {
		return s.boardService.ExcludeChild(ctx, parentId, id)
	}
	return s.laneService.ExcludeChild(ctx, parentId, id)
}

func (s *commandService) updateCard(ctx context.Context, id kernel.ID, name string) error {
	_, err := s.cardService.Name(ctx, id, name)
	return err
}

func (s *commandService) removeCard(ctx context.Context, id kernel.ID, parentId kernel.ID) error {
	err := s.laneService.ExcludeChild(ctx, parentId, id)
	if err == nil {
		err = s.cardService.Remove(ctx, id)
	}
	return err
}

func (s *commandService) layoutBoard(ctx context.Context, id kernel.ID, layout string) error {
	_, err := s.boardService.Layout(ctx, id, layout)
	return err
}

func (s *commandService) getID(key string, payload map[string]string) (kernel.ID, error) {
	value, err := s.getString(key, payload)
	if err != nil {
		return kernel.EmptyID, err
	}

	return kernel.ID(value), nil
}

func (s *commandService) getString(key string, payload map[string]string) (string, error) {
	value, ok := payload[key]
	if !ok {
		s.Errorf("%v is not found in payload\n", key)
		return "", ErrInvalidPayload
	}

	return value, nil
}
