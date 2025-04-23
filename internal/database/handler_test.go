package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/worldline-go/calendar/pkg/models"
	"github.com/worldline-go/query"
	"github.com/worldline-go/test/container/containerpostgres"
	"github.com/worldline-go/types"
)

var migrations = []string{
	"migrations/01_events.sql",
	"migrations/02_relations.sql",
}

type DatabaseSuite struct {
	suite.Suite
	container *containerpostgres.Container
	db        *Database
}

func (s *DatabaseSuite) SetupSuite() {
	s.container = containerpostgres.New(s.T())
	s.container.ExecuteFiles(s.T(), migrations)

	s.db = newDB(s.container.Sqlx(), "public")
}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

func (s *DatabaseSuite) TearDownSuite() {
	s.container.Stop(s.T())
}

func (s *DatabaseSuite) TestAddEvents() {
	events := []models.Event{
		{
			Name:        "New Year",
			Description: "The most wonderful time of the year",
			DateFrom:    types.Time{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
			DateTo:      types.Time{Time: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
		},
	}

	err := s.db.AddEvents(s.T().Context(), events)
	s.Require().NoError(err)

	parse, err := query.Parse("", query.WithExpressionCmp("id", query.ExpressionCmp{
		Operator: query.OperatorEq,
		Field:    "id",
		Value:    events[0].ID,
	}))
	s.Require().NoError(err)

	result, err := s.db.GetEvents(s.T().Context(), parse)
	s.Require().NoError(err)

	s.Require().Len(result, 1)
	s.Require().Equal(events[0].ID, result[0].ID)
	s.Require().Equal(events[0].Name, result[0].Name)
	s.Require().Equal(events[0].Description, result[0].Description)
	s.Require().Equal(events[0].DateFrom.Local(), result[0].DateFrom.Local(), "DateFrom wrong")
	s.Require().Equal(events[0].DateTo.Local(), result[0].DateTo.Local(), "DateTo wrong")
	s.Require().Equal(events[0].RRule, result[0].RRule)
	s.Require().Equal(events[0].Disabled, result[0].Disabled)
	s.Require().Equal(events[0].UpdatedAt.Truncate(time.Millisecond), result[0].UpdatedAt.Truncate(time.Millisecond), "UpdatedAt wrong")
	s.Require().Equal(events[0].UpdatedBy, result[0].UpdatedBy)

	// remove events
	err = s.db.RemoveEvent(s.T().Context(), events[0].ID)
	s.Require().NoError(err)
	// check if removed
	parse, err = query.Parse("", query.WithExpressionCmp("id", query.ExpressionCmp{
		Operator: query.OperatorEq,
		Field:    "id",
		Value:    events[0].ID,
	}))
	s.Require().NoError(err)

	result, err = s.db.GetEvents(s.T().Context(), parse)
	s.Require().NoError(err)
	s.Require().Len(result, 0)
}
