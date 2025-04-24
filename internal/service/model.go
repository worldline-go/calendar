package service

import (
	"context"
	"errors"

	"github.com/worldline-go/query"

	"github.com/worldline-go/calendar/pkg/models"
)

type Database interface {
	AddRelations(ctx context.Context, relations []models.Relation) error
	RemoveRelation(ctx context.Context, q *query.Query) error
	GetRelations(ctx context.Context, q *query.Query) ([]models.Relation, error)
	GetRelationsCount(ctx context.Context, q *query.Query) (uint64, error)

	AddEvents(ctx context.Context, events []models.Event) error
	GetEvents(ctx context.Context, q *query.Query) ([]models.Event, error)
	GetEventsCount(ctx context.Context, q *query.Query) (uint64, error)
	GetEventsWithFunc(ctx context.Context, q *query.Query, fn func(models.Event) error) error
	GetEvent(ctx context.Context, id string) (*models.Event, error)
	UpdateEvent(ctx context.Context, id string, event *models.Event) error
	RemoveEvent(ctx context.Context, id ...string) error
}

var ErrStopLoop = errors.New("stop loop")
