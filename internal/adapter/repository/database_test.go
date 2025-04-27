package repository

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

func (s *DatabaseSuite) TestUpdateEvent() {
	// Add an event
	event := models.Event{
		Name:        "Original Event",
		Description: "Original Description",
		DateFrom:    types.Time{Time: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)},
		DateTo:      types.Time{Time: time.Date(2023, 2, 2, 0, 0, 0, 0, time.UTC)},
	}
	err := s.db.AddEvents(s.T().Context(), []models.Event{event})
	s.Require().NoError(err)

	// Fetch the event to get its ID
	parse, err := query.Parse("", query.WithExpressionCmp("name", query.ExpressionCmp{
		Operator: query.OperatorEq,
		Field:    "name",
		Value:    event.Name,
	}))
	s.Require().NoError(err)
	result, err := s.db.GetEvents(s.T().Context(), parse)
	s.Require().NoError(err)
	s.Require().Len(result, 1)
	eventID := result[0].ID

	// Update the event
	updated := result[0]
	updated.Name = "Updated Event"
	updated.Description = "Updated Description"
	err = s.db.UpdateEvent(s.T().Context(), eventID, &updated)
	s.Require().NoError(err)

	// Fetch and check
	got, err := s.db.GetEvent(s.T().Context(), eventID)
	s.Require().NoError(err)
	s.Require().NotNil(got)
	s.Require().Equal("Updated Event", got.Name)
	s.Require().Equal("Updated Description", got.Description)

	// Cleanup
	_ = s.db.RemoveEvent(s.T().Context(), eventID)
}

func (s *DatabaseSuite) TestGetEventNotFound() {
	got, err := s.db.GetEvent(s.T().Context(), "non-existent-id")
	s.Require().NoError(err)
	s.Require().Nil(got)
}

func (s *DatabaseSuite) TestAddMultipleEvents() {
	events := []models.Event{
		{
			Name:        "Event 1",
			Description: "Desc 1",
			DateFrom:    types.Time{Time: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)},
			DateTo:      types.Time{Time: time.Date(2023, 3, 2, 0, 0, 0, 0, time.UTC)},
		},
		{
			Name:        "Event 2",
			Description: "Desc 2",
			DateFrom:    types.Time{Time: time.Date(2023, 3, 3, 0, 0, 0, 0, time.UTC)},
			DateTo:      types.Time{Time: time.Date(2023, 3, 4, 0, 0, 0, 0, time.UTC)},
		},
	}
	err := s.db.AddEvents(s.T().Context(), events)
	s.Require().NoError(err)

	// Check both events exist
	for _, e := range events {
		parse, err := query.Parse("", query.WithExpressionCmp("name", query.ExpressionCmp{
			Operator: query.OperatorEq,
			Field:    "name",
			Value:    e.Name,
		}))
		s.Require().NoError(err)
		result, err := s.db.GetEvents(s.T().Context(), parse)
		s.Require().NoError(err)
		s.Require().Len(result, 1)
		// Cleanup
		_ = s.db.RemoveEvent(s.T().Context(), result[0].ID)
	}
}

func (s *DatabaseSuite) TestRemoveEventNotFound() {
	// Should not error even if event does not exist
	err := s.db.RemoveEvent(s.T().Context(), "non-existent-id")
	s.Require().NoError(err)
}

func (s *DatabaseSuite) TestAddEventWithAllFields() {
	event := models.Event{
		Name:        "Full Event",
		Description: "All fields set",
		EventGroup:  types.NewNull("group1"),
		DateFrom:    types.Time{Time: time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC)},
		DateTo:      types.Time{Time: time.Date(2023, 4, 2, 0, 0, 0, 0, time.UTC)},
		Tz:          "Europe/Paris",
		AllDay:      true,
		RRule:       "RRULE:FREQ=YEARLY;BYMONTH=4;BYMONTHDAY=1",
		Disabled:    true,
		UpdatedBy:   "tester",
	}
	err := s.db.AddEvents(s.T().Context(), []models.Event{event})
	s.Require().NoError(err)

	// Fetch and check
	parse, err := query.Parse("", query.WithExpressionCmp("name", query.ExpressionCmp{
		Operator: query.OperatorEq,
		Field:    "name",
		Value:    event.Name,
	}))
	s.Require().NoError(err)
	result, err := s.db.GetEvents(s.T().Context(), parse)
	s.Require().NoError(err)
	s.Require().Len(result, 1)
	got := result[0]
	s.Require().Equal(event.Name, got.Name)
	s.Require().Equal(event.Description, got.Description)
	s.Require().Equal(event.EventGroup, got.EventGroup)
	s.Require().Equal(event.Tz, got.Tz)
	s.Require().Equal(event.AllDay, got.AllDay)
	s.Require().Equal(event.RRule, got.RRule)
	s.Require().Equal(event.Disabled, got.Disabled)
	s.Require().Equal(event.UpdatedBy, got.UpdatedBy)

	// Cleanup
	_ = s.db.RemoveEvent(s.T().Context(), got.ID)
}
