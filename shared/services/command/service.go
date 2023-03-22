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
	case kernel.UpdateLaneCommand:
		if name, err := s.getString("name", command.Payload); err != nil {
			result = err
		} else {
			return s.updateLane(ctx, command.ID.WithSet(command.BoardID), name)
		}
	case kernel.RemoveLaneCommand:
		if parentID, err := s.getID("parent_id", command.Payload); err != nil {
			result = err
		} else {
			return s.removeLane(ctx, command.ID.WithSet(command.BoardID), parentID)
		}
	case kernel.UpdateBoardCommand:
		if name, err := s.getString("name", command.Payload); err != nil {
			result = err
		} else {
			return s.updateBoard(ctx, command.BoardID, name)
		}
	case kernel.LayoutBoardCommand:
		if layout, err := s.getString("layout", command.Payload); err != nil {
			result = err
		} else {
			return s.layoutBoard(ctx, command.BoardID, layout)
		}
	case kernel.LayoutLaneCommand:
		if layout, err := s.getString("layout", command.Payload); err != nil {
			result = err
		} else {
			return s.layoutLane(ctx, command.ID.WithSet(command.BoardID), layout)
		}
	case kernel.DescribeBoardCommand:
		if description, err := s.getString("description", command.Payload); err != nil {
			result = err
		} else {
			return s.describeBoard(ctx, command.BoardID, description)
		}
	case kernel.DescribeLaneCommand:
		if description, err := s.getString("description", command.Payload); err != nil {
			result = err
		} else {
			return s.describeLane(ctx, command.ID.WithSet(command.BoardID), description)
		}
	case kernel.DescribeCardCommand:
		if description, err := s.getString("description", command.Payload); err != nil {
			result = err
		} else {
			return s.describeCard(ctx, command.ID.WithSet(command.BoardID), description)
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

func (s *service) describeCard(ctx context.Context, id kernel.MemberID, description string) error {
	return s.cardService.Describe(ctx, id, description)
}

func (s *service) removeCard(ctx context.Context, id kernel.MemberID, parentID kernel.ID) error {
	err := s.laneService.ExcludeChild(ctx, parentID.WithSet(id.SetID), id.ID)
	if err == nil {
		err = s.cardService.Remove(ctx, id)
	}
	return err
}

func (s *service) updateLane(ctx context.Context, id kernel.MemberID, name string) error {
	return s.laneService.Name(ctx, id, name)
}

func (s *service) describeLane(ctx context.Context, id kernel.MemberID, description string) error {
	return s.laneService.Describe(ctx, id, description)
}

func (s *service) removeLane(ctx context.Context, id kernel.MemberID, parentID kernel.ID) error {
	var err error

	if parentID.IsValid() && parentID != id.SetID {
		err = s.laneService.ExcludeChild(ctx, parentID.WithSet(id.SetID), id.ID)
	} else {
		err = s.boardService.ExcludeLane(ctx, id)
	}

	if err == nil {
		err = s.laneService.Remove(ctx, id)
	}

	return err
}

func (s *service) updateBoard(ctx context.Context, id kernel.ID, name string) error {
	return s.boardService.Name(ctx, id, name)
}

func (s *service) layoutBoard(ctx context.Context, id kernel.ID, layout string) error {
	return s.boardService.Layout(ctx, id, layout)
}

func (s *service) describeBoard(ctx context.Context, id kernel.ID, description string) error {
	return s.boardService.Describe(ctx, id, description)
}

func (s *service) layoutLane(ctx context.Context, id kernel.MemberID, layout string) error {
	return s.laneService.Layout(ctx, id, layout)
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
