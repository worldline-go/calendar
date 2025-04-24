package ports

import (
	"context"
	"errors"

	"github.com/worldline-go/calendar/internal/database"
	"github.com/worldline-go/calendar/internal/service"
	"github.com/worldline-go/calendar/pkg/models"
	"github.com/worldline-go/query"
)

type PortDB struct {
	db *database.Database
}

var _ service.Database = (*PortDB)(nil)

func NewDatabase(db *database.Database) *PortDB {
	return &PortDB{
		db: db,
	}
}

func (p *PortDB) AddRelations(ctx context.Context, relations []models.Relation) error {
	return p.db.AddRelations(ctx, relations)
}

func (p *PortDB) RemoveRelation(ctx context.Context, q *query.Query) error {
	return p.db.RemoveRelation(ctx, q)
}

func (p *PortDB) GetRelations(ctx context.Context, q *query.Query) ([]models.Relation, error) {
	return p.db.GetRelations(ctx, q)
}

func (p *PortDB) GetRelationsCount(ctx context.Context, q *query.Query) (uint64, error) {
	return p.db.GetRelationsCount(ctx, q)
}

func (p *PortDB) AddEvents(ctx context.Context, events []models.Event) error {
	return p.db.AddEvents(ctx, events)
}

func (p *PortDB) GetEvents(ctx context.Context, q *query.Query) ([]models.Event, error) {
	return p.db.GetEvents(ctx, q)
}

func (p *PortDB) GetEventsCount(ctx context.Context, q *query.Query) (uint64, error) {
	return p.db.GetEventsCount(ctx, q)
}

func (p *PortDB) GetEventsWithFunc(ctx context.Context, q *query.Query, fn func(models.Event) error) error {
	return p.db.GetEventsWithFunc(ctx, q, func(e models.Event) error {
		if err := fn(e); err != nil {
			if errors.Is(err, service.ErrStopLoop) {
				return database.ErrStopLoop
			}

			return err
		}

		return nil
	})
}

func (p *PortDB) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	return p.db.GetEvent(ctx, id)
}

func (p *PortDB) UpdateEvent(ctx context.Context, id string, event *models.Event) error {
	return p.db.UpdateEvent(ctx, id, event)
}

func (p *PortDB) RemoveEvent(ctx context.Context, id ...string) error {
	return p.db.RemoveEvent(ctx, id...)
}
