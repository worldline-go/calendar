package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/worldline-go/calendar/pkg/models"
	"github.com/worldline-go/query"
	"github.com/worldline-go/test/container"
)

var migrations = []string{
	"../../migrations/01_calendar.sql",
	"../../migrations/02_relation.sql",
}

type DatabaseSuite struct {
	suite.Suite
	container *container.PostgresContainer
	db        *Database
}

func (s *DatabaseSuite) SetupSuite() {
	s.container = container.Postgres(s.T())

	s.container.CreateSchema(s.T(), "public")
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
	calendar := &models.Event{
		Name:        "New Year",
		Description: "The most wonderful time of the year",
	}

	err := s.db.AddEvents(s.T().Context(), calendar)
	s.Require().NoError(err)

	parse, err := query.Parse("", query.WithExpressionCmp("id", query.ExpressionCmp{
		Operator: query.OperatorEq,
		Field:    "id",
		Value:    calendar.ID,
	}))
	s.Require().NoError(err)

	result, err := s.db.GetEvents(s.T().Context(), parse)
	s.Require().NoError(err)

	s.Require().Len(result, 1)
	s.Require().Equal(calendar.ID, result[0].ID)
	s.Require().Equal(calendar.Name, result[0].Name)
	s.Require().Equal(calendar.Description, result[0].Description)
	s.Require().Equal(calendar.DateFrom, result[0].DateFrom)
	s.Require().Equal(calendar.DateTo, result[0].DateTo)
	s.Require().Equal(calendar.RRule, result[0].RRule)
	s.Require().Equal(calendar.Disabled, result[0].Disabled)
	s.Require().Equal(calendar.UpdatedAt.Truncate(time.Millisecond), result[0].UpdatedAt.Truncate(time.Millisecond))
	s.Require().Equal(calendar.UpdatedBy, result[0].UpdatedBy)
}
