package port

import (
	"context"
	"io"
	"time"

	"github.com/worldline-go/query"
	"github.com/worldline-go/types"

	"github.com/worldline-go/calendar/internal/core/domain"
)

type CalendarPort interface {
	AddRelations(ctx context.Context, relations []domain.Relation) error
	RemoveRelation(ctx context.Context, q *query.Query) error
	GetRelations(ctx context.Context, q *query.Query) ([]domain.Relation, error)
	GetRelationsCount(ctx context.Context, q *query.Query) (uint64, error)
	AddEvents(ctx context.Context, events []domain.Event) error
	GetEvents(ctx context.Context, q *query.Query) ([]domain.Event, error)
	GetEventsCount(ctx context.Context, q *query.Query) (uint64, error)
	GetEventsWithFunc(ctx context.Context, q *query.Query, fn func(domain.Event) error) error
	GetEvent(ctx context.Context, id string) (*domain.Event, error)
	UpdateEvent(ctx context.Context, id string, event *domain.Event) error
	RemoveEvent(ctx context.Context, id ...string) error
}

type CalendarService interface {
	AddRelations(ctx context.Context, relations []domain.Relation) error
	RemoveRelation(ctx context.Context, q *query.Query) error
	GetRelations(ctx context.Context, q *query.Query) ([]domain.Relation, error)
	GetRelationsCount(ctx context.Context, q *query.Query) (uint64, error)
	AddEvents(ctx context.Context, events []domain.Event) error
	GetEvents(ctx context.Context, q *query.Query) ([]domain.Event, error)
	GetEventsCount(ctx context.Context, q *query.Query) (uint64, error)
	GetEvent(ctx context.Context, id string) (*domain.Event, error)
	UpdateEvent(ctx context.Context, id string, event *domain.Event) error
	RemoveEvent(ctx context.Context, id ...string) error

	AddIcal(ctx context.Context, data io.Reader, tz *time.Location, group types.Null[string], updatedBy string) error
	GetEventsICS(ctx context.Context, q *query.Query) ([]domain.Event, error)
}
