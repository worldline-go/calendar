package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/worldline-go/query"
	"github.com/worldline-go/types"

	"github.com/worldline-go/calendar/pkg/models"
)

type Service struct {
	Database Database
}

func New(db Database) *Service {
	return &Service{
		Database: db,
	}
}

// WorkDay returns the next workday after the given date.
func (s *Service) WorkDay(ctx context.Context, date types.Time) (types.Time, error) {
	return types.Time{Time: time.Now()}, nil
}

// //////////////////////////////////////////////////////////////
// Database
// //////////////////////////////////////////////////////////////

func (s *Service) AddHoliday(ctx context.Context, holiday ...*models.Holiday) error {
	if err := s.Database.AddHoliday(ctx, holiday...); err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveHoliday(ctx context.Context, id string) error {
	err := s.Database.RemoveHoliday(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetHolidaysCount(ctx context.Context, q *query.Query) (int64, error) {
	count, err := s.Database.GetHolidaysCount(ctx, q)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Service) GetHolidays(ctx context.Context, q *query.Query) ([]*models.Holiday, error) {
	if q.HasAny("date") {
		var holidays []*models.Holiday

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

		yearCheck := false
		year := int64(0)
		if qYear, _ := q.Values["year"]; len(qYear) > 0 {
			qYearStr, ok := qYear[0].Value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid year format")
			}

			var err error
			year, err = strconv.ParseInt(qYearStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid year format: %w", err)
			}

			yearCheck = true
		}

		limitOrg := q.CloneLimit()
		offsetOrg := q.CloneOffset()

		limit := q.CloneLimit()
		offset := q.CloneOffset()

		q.Limit = nil
		q.Offset = nil

		defer func() {
			q.Limit = limitOrg
			q.Offset = offsetOrg
		}()

		err := s.Database.GetHolidaysWithFunc(ctx, q, func(h *models.Holiday) error {
			if offset != nil {
				if *offset > 0 {
					*offset--

					return nil
				}
			}

			if dateCheck {
				ok, err := CheckDate(qDateCheck, h.DateFrom, h.DateTo, h.RRule)
				if err != nil {
					return err
				}

				if !ok {
					return nil
				}
			}

			if yearCheck {
				ok, err := CheckYear(int(year), h.DateFrom, h.DateTo, h.RRule)
				if err != nil {
					return err
				}

				if !ok {
					return nil
				}
			}

			modifyYear := 0
			if yearCheck {
				modifyYear = int(year)
			} else {
				modifyYear = qDateCheck.Year()
			}

			// modify the date's year
			if h.DateFrom.Valid {
				h.DateFrom.V = types.Time{Time: ChangeYear(h.DateFrom.V.Time, modifyYear)}
			}

			if h.DateTo.Valid {
				h.DateTo.V = types.Time{Time: ChangeYear(h.DateTo.V.Time, modifyYear)}
			}

			holidays = append(holidays, h)

			if limit != nil {
				if *limit > 0 {
					*limit--
				}

				if *limit == 0 {
					return models.ErrStopLoop
				}
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		return holidays, nil
	}

	holidays, err := s.Database.GetHolidays(ctx, q)
	if err != nil {
		return nil, err
	}

	return holidays, nil
}

func (s *Service) GetHoliday(ctx context.Context, id string) (*models.Holiday, error) {
	holiday, err := s.Database.GetHoliday(ctx, id)
	if err != nil {
		return nil, err
	}

	return holiday, nil
}

func (s *Service) UpdateHoliday(ctx context.Context, id string, holiday *models.Holiday) error {
	err := s.Database.UpdateHoliday(ctx, id, holiday)
	if err != nil {
		return err
	}

	return nil
}

// ///////////////////////////////////////////////////////////////
// Relations
// //////////////////////////////////////////////////////////////

func (s *Service) AddRelation(ctx context.Context, relations ...*models.Relation) error {
	if err := s.Database.AddRelation(ctx, relations...); err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveRelation(ctx context.Context, id string) error {
	err := s.Database.RemoveRelation(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetRelation(ctx context.Context, id string) (*models.Relation, error) {
	relation, err := s.Database.GetRelation(ctx, id)
	if err != nil {
		return nil, err
	}

	return relation, nil
}

func (s *Service) GetRelations(ctx context.Context, q *query.Query) ([]*models.Relation, error) {
	relations, err := s.Database.GetRelations(ctx, q)
	if err != nil {
		return nil, err
	}

	return relations, nil
}

func (s *Service) GetRelationsCount(ctx context.Context, q *query.Query) (int64, error) {
	count, err := s.Database.GetRelationsCount(ctx, q)
	if err != nil {
		return 0, err
	}

	return count, nil
}
