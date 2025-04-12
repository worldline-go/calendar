package service

import (
	"context"

	"github.com/worldline-go/query"

	"github.com/worldline-go/calendar/pkg/models"
)

type Database interface {
	AddRelation(ctx context.Context, relations ...*models.Relation) error
	RemoveRelation(ctx context.Context, id string) error
	GetRelation(ctx context.Context, id string) (*models.Relation, error)
	GetRelations(ctx context.Context, q *query.Query) ([]*models.Relation, error)
	GetRelationsCount(ctx context.Context, q *query.Query) (int64, error)
	GetHolidaysWithFunc(ctx context.Context, q *query.Query, fn func(*models.Holiday) error) error

	AddHoliday(ctx context.Context, holidays ...*models.Holiday) error
	GetHolidays(ctx context.Context, q *query.Query) ([]*models.Holiday, error)
	GetHolidaysCount(ctx context.Context, q *query.Query) (int64, error)
	GetHoliday(ctx context.Context, id string) (*models.Holiday, error)
	UpdateHoliday(ctx context.Context, id string, holiday *models.Holiday) error
	RemoveHoliday(ctx context.Context, id string) error
}
