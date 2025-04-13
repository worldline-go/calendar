package database

import (
	"context"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/oklog/ulid/v2"
	"github.com/worldline-go/query"
	"github.com/worldline-go/query/adapter/adaptergoqu"
	"github.com/worldline-go/types"

	"github.com/worldline-go/calendar/pkg/models"
)

type QeuryHoliday struct {
	ID       *int64
	Provider *int64
	Country  *string
}

var (
	TableEventsStr    = "calendar_events"
	TableRelationsStr = "calendar_relations"

	Schema        exp.IdentifierExpression
	TableEvents   exp.AliasedExpression
	TableRelation exp.AliasedExpression
)

func setSchema(schema string) {
	Schema = goqu.S(schema)
	TableEvents = Schema.Table(TableEventsStr).As(TableEventsStr)
	TableRelation = Schema.Table(TableRelationsStr).As(TableRelationsStr)
}

func (db *Database) AddEvents(ctx context.Context, events ...*models.Event) error {
	updatedAt := types.Time{Time: time.Now()}

	for i := range events {
		events[i].ID = ulid.Make().String()
		events[i].UpdatedAt = updatedAt
	}

	var holidayResult []*models.Event

	err := db.q.Insert(TableEvents).
		Rows(events).
		Executor().ScanStructsContext(ctx, &holidayResult)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) getEventsSelect(q *query.Query) *goqu.SelectDataset {
	selectDataSet := adaptergoqu.Select(q, db.q.From(TableEvents),
		adaptergoqu.WithDefaultSelect(
			TableEventsStr+".id",
			TableEventsStr+".name",
			TableEventsStr+".description",
			TableEventsStr+".date_from",
			TableEventsStr+".date_to",
			TableEventsStr+".years",
			TableEventsStr+".disabled",
			TableEventsStr+".updated_at",
			TableEventsStr+".updated_by",
		),
		adaptergoqu.WithRename(map[string]string{
			"code":    TableRelationsStr + ".code",
			"country": TableRelationsStr + ".country",
		}),
	)

	if q.HasAny("code", "country") {
		selectDataSet = selectDataSet.RightJoin(TableRelation, goqu.On(goqu.Ex{TableRelationsStr + ".holiday_id": goqu.I("holiday.id")}))
	}

	return selectDataSet
}

func (db *Database) GetEventsCount(ctx context.Context, q *query.Query) (int64, error) {
	count, err := db.getEventsSelect(q).CountContext(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (db *Database) GetEvents(ctx context.Context, q *query.Query) ([]*models.Event, error) {
	var events []*models.Event

	if err := db.getEventsSelect(q).Executor().ScanStructsContext(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (db *Database) GetEventsWithFunc(ctx context.Context, q *query.Query, fn func(*models.Event) error) error {
	scanner, err := db.getEventsSelect(q).Executor().ScannerContext(ctx)
	if err != nil {
		return err
	}
	defer scanner.Close()

	for scanner.Next() {
		var event models.Event
		if err := scanner.ScanStruct(&event); err != nil {
			return err
		}
		if err := fn(&event); err != nil {
			if errors.Is(err, models.ErrStopLoop) {
				break
			}

			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (db *Database) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	var event models.Event

	found, err := db.q.From(TableEvents).
		Where(goqu.Ex{
			"id": id,
		}).
		Executor().ScanStructContext(ctx, &event)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &event, nil
}

func (db *Database) UpdateEvent(ctx context.Context, id string, event *models.Event) error {
	_, err := db.q.Update(TableEvents).
		Set(event).
		Where(goqu.Ex{
			"id": id,
		}).
		Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveEvent(ctx context.Context, id string) error {
	_, err := db.q.Delete(TableEvents).
		Where(goqu.Ex{
			"id": id,
		}).
		Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

// /////////////////////////////////////////////////////////////
// Relation
// /////////////////////////////////////////////////////////////

func (db *Database) AddRelations(ctx context.Context, relations ...*models.Relation) error {
	updatedAt := types.Time{Time: time.Now()}

	for i := range relations {
		relations[i].ID = ulid.Make().String()
		relations[i].UpdatedAt = updatedAt
	}

	_, err := db.q.Insert(TableRelation).
		Rows(relations).
		Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveRelation(ctx context.Context, id string) error {
	_, err := db.q.Delete(TableRelation).
		Where(goqu.Ex{
			"id": id,
		}).
		Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetRelationsCount(ctx context.Context, q *query.Query) (int64, error) {
	count, err := adaptergoqu.Select(q, db.q.From(TableRelation)).CountContext(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (db *Database) GetRelations(ctx context.Context, q *query.Query) ([]*models.Relation, error) {
	var relations []*models.Relation

	if err := adaptergoqu.Select(q, db.q.From(TableRelation)).Executor().ScanStructsContext(ctx, &relations); err != nil {
		return nil, err
	}

	return relations, nil
}

func (db *Database) GetRelation(ctx context.Context, id string) (*models.Relation, error) {
	var relation models.Relation

	found, err := db.q.From(TableRelation).
		Where(goqu.Ex{
			"id": id,
		}).
		Executor().ScanStructContext(ctx, &relation)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &relation, nil
}
