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
	TableHolidayStr  = "holiday"
	TableRelationStr = "holiday_relation"

	Schema        exp.IdentifierExpression
	TableHoliday  exp.AliasedExpression
	TableRelation exp.AliasedExpression
)

func setSchema(schema string) {
	Schema = goqu.S(schema)
	TableHoliday = Schema.Table(TableHolidayStr).As(TableHolidayStr)
	TableRelation = Schema.Table(TableRelationStr).As(TableRelationStr)
}

func (db *Database) AddHoliday(ctx context.Context, holidays ...*models.Holiday) error {
	updatedAt := types.Time{Time: time.Now()}

	for i := range holidays {
		holidays[i].ID = ulid.Make().String()
		holidays[i].UpdatedAt = updatedAt
	}

	var holidayResult []*models.Holiday

	err := db.q.Insert(TableHoliday).
		Rows(holidays).
		Executor().ScanStructsContext(ctx, &holidayResult)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) getHolidaysSelect(q *query.Query) *goqu.SelectDataset {
	selectDataSet := adaptergoqu.Select(q, db.q.From(TableHoliday),
		adaptergoqu.WithDefaultSelect(
			TableHolidayStr+".id",
			TableHolidayStr+".name",
			TableHolidayStr+".description",
			TableHolidayStr+".date_from",
			TableHolidayStr+".date_to",
			TableHolidayStr+".years",
			TableHolidayStr+".disabled",
			TableHolidayStr+".updated_at",
			TableHolidayStr+".updated_by",
		),
		adaptergoqu.WithRename(map[string]string{
			"code":    TableRelationStr + ".code",
			"country": TableRelationStr + ".country",
		}),
	)

	if q.HasAny("code", "country") {
		selectDataSet = selectDataSet.RightJoin(TableRelation, goqu.On(goqu.Ex{TableRelationStr + ".holiday_id": goqu.I("holiday.id")}))
	}

	return selectDataSet
}

func (db *Database) GetHolidaysCount(ctx context.Context, q *query.Query) (int64, error) {
	count, err := db.getHolidaysSelect(q).CountContext(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (db *Database) GetHolidays(ctx context.Context, q *query.Query) ([]*models.Holiday, error) {
	var holidays []*models.Holiday

	if err := db.getHolidaysSelect(q).Executor().ScanStructsContext(ctx, &holidays); err != nil {
		return nil, err
	}

	return holidays, nil
}

func (db *Database) GetHolidaysWithFunc(ctx context.Context, q *query.Query, fn func(*models.Holiday) error) error {
	scanner, err := db.getHolidaysSelect(q).Executor().ScannerContext(ctx)
	if err != nil {
		return err
	}
	defer scanner.Close()

	for scanner.Next() {
		var holiday models.Holiday
		if err := scanner.ScanStruct(&holiday); err != nil {
			return err
		}
		if err := fn(&holiday); err != nil {
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

func (db *Database) GetHoliday(ctx context.Context, id string) (*models.Holiday, error) {
	var holiday models.Holiday

	found, err := db.q.From(TableHoliday).
		Where(goqu.Ex{
			"id": id,
		}).
		Executor().ScanStructContext(ctx, &holiday)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &holiday, nil
}

func (db *Database) UpdateHoliday(ctx context.Context, id string, holiday *models.Holiday) error {
	_, err := db.q.Update(TableHoliday).
		Set(holiday).
		Where(goqu.Ex{
			"id": id,
		}).
		Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveHoliday(ctx context.Context, id string) error {
	_, err := db.q.Delete(TableHoliday).
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

func (db *Database) AddRelation(ctx context.Context, relations ...*models.Relation) error {
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
