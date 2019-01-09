package services

import (
	"context"
	"errors"
	"strconv"

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
	Execute(context.Context, kernel.Id, kernel.CommandType, map[string]string) error
}

type commandService struct {
	logger.Logger
	boardService BoardService
	laneService  LaneService
}

func (s *commandService) Execute(ctx context.Context, id kernel.Id, command kernel.CommandType, payload map[string]string) error {
	switch command {
	case kernel.InsertBefore: //todo
	case kernel.InsertAfter: //todo
	case kernel.AppendChild:
		if laneID, err := s.getID("lane_id", payload); err != nil {
			return err
		} else {
			return s.appendChild(ctx, id, laneID)
		}
	case kernel.ExcludeChild:
		if laneID, err := s.getID("lane_id", payload); err != nil {
			return err
		} else {
			return s.excludeChild(ctx, id, laneID)
		}
	case kernel.UpdateCard: //todo
	case kernel.RemoveCard: //todo
	case kernel.LayoutBoard:
		if layout, err := s.getString("layout", payload); err != nil {
			return err
		} else {
			return s.layoutBoard(ctx, id, layout)
		}
	}

	return ErrInvalidCommandType
}

func (s *commandService) insertBefore(ctx context.Context, id kernel.Id, relativeId kernel.Id) error {
	return nil
}

func (s *commandService) InsertAfter(ctx context.Context, id kernel.Id, relativeId kernel.Id) error {
	return nil
}

func (s *commandService) appendChild(ctx context.Context, id kernel.Id, parentId kernel.Id) error {
	return s.laneService.AppendChild(ctx, parentId, id)
}

func (s *commandService) excludeChild(ctx context.Context, id kernel.Id, parentId kernel.Id) error {
	return s.laneService.ExcludeChild(ctx, parentId, id)
}

func (s *commandService) UpdateCard(ctx context.Context, id kernel.Id, relativeId kernel.Id) error {
	return nil
}

func (s *commandService) RemoveCard(ctx context.Context, id kernel.Id) error {
	return nil
}

func (s *commandService) layoutBoard(ctx context.Context, id kernel.Id, layout string) error {
	_, err := s.boardService.Layout(ctx, id, layout)
	return err
}

func (s *commandService) getID(key string, payload map[string]string) (kernel.Id, error) {
	value, err := s.getString(key, payload)
	if err != nil {
		return kernel.Id(0), err
	}

	id, err := strconv.Atoi(value)
	if err != nil {
		s.Errorln(err)
		return "", ErrInvalidPayload
	}

	return kernel.Id(id), nil
}

func (s *commandService) getString(key string, payload map[string]string) (string, error) {
	value, ok := payload[key]
	if !ok {
		s.Errorf("%v is not found in payload\n", key)
		return "", ErrInvalidPayload
	}

	return value, nil
}
