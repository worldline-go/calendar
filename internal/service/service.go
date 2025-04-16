package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/worldline-go/cache"
	"github.com/worldline-go/cache/store/memory"
	"github.com/worldline-go/query"
	"github.com/worldline-go/types"

	"github.com/worldline-go/calendar/pkg/ics"
	"github.com/worldline-go/calendar/pkg/models"
)

type Service struct {
	db    Database
	cache cache.Cacher[string, *ics.RRule]
	m     sync.RWMutex
}

func New(ctx context.Context, db Database) (*Service, error) {
	cache, err := cache.New[string, *ics.RRule](ctx,
		memory.Store,
		cache.WithStoreConfig(memory.Config{
			MaxItems: 200,
			TTL:      30 * time.Minute,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	return &Service{
		cache: cache,
		db:    db,
	}, nil
}

// WorkDay returns the next workday after the given date.
func (s *Service) WorkDay(ctx context.Context, date types.Time) (types.Time, error) {
	return types.Time{Time: time.Now()}, nil
}

// //////////////////////////////////////////////////////////////
// Database
// //////////////////////////////////////////////////////////////

func (s *Service) AddEvents(ctx context.Context, events ...*models.Event) error {
	if err := s.db.AddEvents(ctx, events...); err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveEvent(ctx context.Context, id string) error {
	err := s.db.RemoveEvent(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetEventsCount(ctx context.Context, q *query.Query) (uint64, error) {
	count, err := s.db.GetEventsCount(ctx, q)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Service) GetEvents(ctx context.Context, q *query.Query) ([]*models.Event, error) {
	if q.HasAny("date") {
		var events []*models.Event

		dateCheck := false
		qDateCheck := types.Time{}
		if qDate, _ := q.Values["date"]; len(qDate) > 0 {
			qDateStr, ok := qDate[0].Value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid date format")
			}
			if err := qDateCheck.Parse(qDateStr); err != nil {
				return nil, fmt.Errorf("invalid date format: %w", err)
			}

			dateCheck = true
		}

		err := s.db.GetEventsWithFunc(ctx, q, func(h *models.Event) error {
			if h.Disabled {
				return nil
			}

			icsRule, err := s.getRRule(ctx, h.RRule)
			if err != nil {
				return fmt.Errorf("failed to get rrule: %w", err)
			}

			if dateCheck {
				start, stop, ok := ics.MatchRRuleAt(icsRule, h.DateFrom.Time, h.DateTo.Time, qDateCheck.Time)
				if !ok {
					return nil
				}

				h.DateFrom.Time = start
				h.DateTo.Time = stop
			}

			events = append(events, h)

			return nil
		})
		if err != nil {
			return nil, err
		}

		return events, nil
	}

	holidays, err := s.db.GetEvents(ctx, q)
	if err != nil {
		return nil, err
	}

	return holidays, nil
}

func (s *Service) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	holiday, err := s.db.GetEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	return holiday, nil
}

func (s *Service) UpdateEvent(ctx context.Context, id string, event *models.Event) error {
	err := s.db.UpdateEvent(ctx, id, event)
	if err != nil {
		return err
	}

	return nil
}

// ///////////////////////////////////////////////////////////////
// Relations
// //////////////////////////////////////////////////////////////

func (s *Service) AddRelations(ctx context.Context, relations ...*models.Relation) error {
	if err := s.db.AddRelations(ctx, relations...); err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveRelation(ctx context.Context, id string) error {
	err := s.db.RemoveRelation(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetRelation(ctx context.Context, id string) (*models.Relation, error) {
	relation, err := s.db.GetRelation(ctx, id)
	if err != nil {
		return nil, err
	}

	return relation, nil
}

func (s *Service) GetRelations(ctx context.Context, q *query.Query) ([]*models.Relation, error) {
	relations, err := s.db.GetRelations(ctx, q)
	if err != nil {
		return nil, err
	}

	return relations, nil
}

func (s *Service) GetRelationsCount(ctx context.Context, q *query.Query) (int64, error) {
	count, err := s.db.GetRelationsCount(ctx, q)
	if err != nil {
		return 0, err
	}

	return count, nil
}
