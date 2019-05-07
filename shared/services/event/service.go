package event

import (
	"context"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/services/notification"
	"github.com/dmibod/kanban/shared/services/persistence"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type service struct {
	logger.Logger
	notificationService notification.Service
	persistenceService  persistence.Service
}

// CreateService instance
func CreateService(n notification.Service, p persistence.Service, l logger.Logger) Service {
	return &service{
		notificationService: n,
		persistenceService:  p,
		Logger:              l,
	}
}

func (s *service) Execute(c context.Context, f func(event.Bus) error) error {
	return event.Execute(func(bus event.Bus) error {
		s.persistenceService.Listen(bus)
		s.notificationService.Listen(bus)
		if err := f(bus); err != nil {
			s.Errorln(err)
			return err
		}
		return bus.Fire(c)
	})
}
