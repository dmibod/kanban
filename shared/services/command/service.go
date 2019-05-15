package command

import (
	"context"

	"github.com/dmibod/kanban/shared/services/board"
	"github.com/dmibod/kanban/shared/services/card"
	"github.com/dmibod/kanban/shared/services/lane"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type service struct {
	logger.Logger
	boardService board.Service
	laneService  lane.Service
	cardService  card.Service
}

// CreateService instance
func CreateService(bs board.Service, ls lane.Service, cs card.Service, l logger.Logger) Service {
	return &service{
		Logger:       l,
		boardService: bs,
		laneService:  ls,
		cardService:  cs,
	}
}

func (s *service) Execute(ctx context.Context, command kernel.Command) error {
	result := ErrInvalidCommandType

	switch command.Type {
	case kernel.InsertBeforeCommand: //todo
	case kernel.InsertAfterCommand: //todo
	case kernel.AppendChildCommand:
		if parentID, err := s.getID("parent_id", command.Payload); err != nil {
			result = err
		} else {
			return s.appendChild(ctx, command.ID.WithSet(command.BoardID), parentID)
		}
	case kernel.ExcludeChildCommand:
		if parentID, err := s.getID("parent_id", command.Payload); err != nil {
			result = err
		} else {
			return s.excludeChild(ctx, command.ID.WithSet(command.BoardID), parentID)
		}
	case kernel.UpdateCardCommand:
		if name, err := s.getString("name", command.Payload); err != nil {
			result = err
		} else {
			return s.updateCard(ctx, command.ID.WithSet(command.BoardID), name)
		}
	case kernel.RemoveCardCommand:
		if parentID, err := s.getID("parent_id", command.Payload); err != nil {
			result = err
		} else {
			return s.removeCard(ctx, command.ID.WithSet(command.BoardID), parentID)
		}
	case kernel.LayoutBoardCommand:
		if layout, err := s.getString("layout", command.Payload); err != nil {
			result = err
		} else {
			return s.layoutBoard(ctx, command.BoardID, layout)
		}
	}

	return result
}

func (s *service) insertBefore(ctx context.Context, id kernel.ID, relativeID kernel.ID) error {
	return nil
}

func (s *service) insertAfter(ctx context.Context, id kernel.ID, relativeID kernel.ID) error {
	return nil
}

func (s *service) appendChild(ctx context.Context, id kernel.MemberID, parentID kernel.ID) error {
	if _, err := s.boardService.GetByID(ctx, parentID); err == nil {
		return s.boardService.AppendLane(ctx, id)
	}
	return s.laneService.AppendChild(ctx, parentID.WithSet(id.SetID), id.ID)
}

func (s *service) excludeChild(ctx context.Context, id kernel.MemberID, parentID kernel.ID) error {
	if _, err := s.boardService.GetByID(ctx, parentID); err == nil {
		return s.boardService.ExcludeLane(ctx, id)
	}
	return s.laneService.ExcludeChild(ctx, parentID.WithSet(id.SetID), id.ID)
}

func (s *service) updateCard(ctx context.Context, id kernel.MemberID, name string) error {
	return s.cardService.Name(ctx, id, name)
}

func (s *service) removeCard(ctx context.Context, id kernel.MemberID, parentID kernel.ID) error {
	err := s.laneService.ExcludeChild(ctx, parentID.WithSet(id.SetID), id.ID)
	if err == nil {
		err = s.cardService.Remove(ctx, id)
	}
	return err
}

func (s *service) layoutBoard(ctx context.Context, id kernel.ID, layout string) error {
	return s.boardService.Layout(ctx, id, layout)
}

func (s *service) getID(key string, payload map[string]string) (kernel.ID, error) {
	value, err := s.getString(key, payload)
	if err != nil {
		return kernel.EmptyID, err
	}

	return kernel.ID(value), nil
}

func (s *service) getString(key string, payload map[string]string) (string, error) {
	value, ok := payload[key]
	if !ok {
		s.Errorf("%v is not found in payload\n", key)
		return "", ErrInvalidPayload
	}

	return value, nil
}
