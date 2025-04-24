package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/worldline-go/query"
	"github.com/worldline-go/rest"
	"github.com/worldline-go/rest/server"
	"github.com/worldline-go/types"

	"github.com/worldline-go/calendar/internal/service"
	"github.com/worldline-go/calendar/pkg/ical"
	"github.com/worldline-go/calendar/pkg/models"
)

type HTTP struct {
	Service *service.Service

	Validator QueryValidator
}

type QueryValidator struct {
	GetEvents    *query.Validator
	DeleteEvents *query.Validator

	GetRelations    *query.Validator
	DeleteRelations *query.Validator

	GetEventsDate *query.Validator
	GetICS        *query.Validator
}

var DefaultLimit uint64 = 25

func NewHTTP(svc *service.Service) (*HTTP, error) {
	validatorGetEvents, err := query.NewValidator(
		query.WithField(query.WithNotAllowed()),
		query.WithSort(query.WithIn("id", "entity", "event_group", "name", "description", "disabled", "date_from", "date_to", "updated_at", "updated_by")),
		query.WithValues(query.WithIn("id", "entity", "event_group", "name", "description", "disabled", "updated_by")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator for GetEvents: %w", err)
	}

	validatorDeleteEvents, err := query.NewValidator(
		query.WithField(query.WithNotAllowed()),
		query.WithValues(query.WithIn("id")),
		query.WithValue("id", query.WithOperator(query.OperatorEq, query.OperatorIn), query.WithNotEmpty()),
		query.WithLimit(query.WithNotAllowed()),
		query.WithOffset(query.WithNotAllowed()),
		query.WithSort(query.WithNotAllowed()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator for DeleteEvents: %w", err)
	}

	validatorGetRelations, err := query.NewValidator(
		query.WithField(query.WithNotAllowed()),
		query.WithSort(query.WithIn("event_id", "event_group", "entity")),
		query.WithValues(query.WithIn("entity", "event_id", "event_group")),
		query.WithValue("entity", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithValue("event_id", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithValue("event_group", query.WithOperator(query.OperatorEq, query.OperatorIn)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator for GetRelations: %w", err)
	}

	validatorDeleteRelations, err := query.NewValidator(
		query.WithField(query.WithNotAllowed()),
		query.WithValues(query.WithIn("entity", "event_id", "event_group")),
		query.WithValue("entity", query.WithOperator(query.OperatorEq, query.OperatorIn), query.WithNotEmpty()),
		query.WithValue("event_id", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithValue("event_group", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithLimit(query.WithNotAllowed()),
		query.WithOffset(query.WithNotAllowed()),
		query.WithSort(query.WithNotAllowed()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator for DeleteRelations: %w", err)
	}

	validatorGetEventsDate, err := query.NewValidator(
		query.WithField(query.WithNotAllowed()),
		query.WithValues(query.WithIn("entity", "event_group", "date")),
		query.WithValue("entity", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithValue("event_group", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithValue("date", query.WithOperator(query.OperatorEq), query.WithNotEmpty()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator for GetEventsDate: %w", err)
	}

	validatorGetICS, err := query.NewValidator(
		query.WithField(query.WithNotAllowed()),
		query.WithValues(query.WithIn("entity", "event_group", "year")),
		query.WithValue("entity", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithValue("event_group", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithValue("year", query.WithOperator(query.OperatorEq, query.OperatorIn)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator for GetICS: %w", err)
	}

	return &HTTP{
		Service: svc,
		Validator: QueryValidator{
			GetEvents:       validatorGetEvents,
			DeleteEvents:    validatorDeleteEvents,
			DeleteRelations: validatorDeleteRelations,
			GetRelations:    validatorGetRelations,
			GetEventsDate:   validatorGetEventsDate,
			GetICS:          validatorGetICS,
		},
	}, nil
}

func (h *HTTP) RegisterRoutes(g *echo.Group) {
	g.GET("/events", h.GetEvents)
	g.POST("/events", h.AddEvents)
	g.DELETE("/events", h.DeleteEvents)

	g.GET("/events/{id}", h.GetEvent)
	g.DELETE("/events/{id}", h.DeleteEvent)
	g.PUT("/events/{id}", h.PutEvent)

	g.GET("/relations", h.GetRelations)
	g.POST("/relations", h.AddRelations)
	g.DELETE("/relations", h.DeleteRelations)

	g.GET("/holidays", h.Holidays)
	g.POST("/ics", h.AddICS)
	g.GET("/ics", h.GetICS)
}

// @Summary GetEvents
// @Description GetEvents
// @Param id query string false "id"
// @Param name query string false "name"
// @Param description query string false "description"
// @Param event_group query string false "event_group"
// @Param entity query string false "entity for relation"
// @Param disabled query bool false "disabled"
// @Param limit query int false "limit" default(25)
// @Param offset query int false "offset"
// @Success 200 {object} rest.Response[[]models.Event]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /events [get]
// @Tags Events
func (h *HTTP) GetEvents(c echo.Context) error {
	q, err := query.ParseWithValidator(
		c.QueryString(),
		h.Validator.GetEvents,
		query.WithDefaultLimit(DefaultLimit),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	events, err := h.Service.GetEvents(c.Request().Context(), q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if len(events) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "no events found")
	}

	count, err := h.Service.GetEventsCount(c.Request().Context(), q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "events count failed").SetInternal(err)
	}

	return c.JSON(http.StatusOK, rest.Response[[]models.Event]{
		Meta: &rest.Meta{
			TotalItemCount: count,
			Limit:          q.GetLimit(),
			Offset:         q.GetOffset(),
		},
		Payload: events,
	})
}

// @Summary AddEvents
// @Description AddEvents
// @Param body body []models.Event true "Event"
// @Success 200 {object} rest.Response[[]string]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /events [post]
// @Tags Events
func (h *HTTP) AddEvents(c echo.Context) error {
	v := []models.Event{}
	if err := rest.BindJSONList(c.Request().Body, &v); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	updatedBy := server.GetUser(c)
	for i := range v {
		v[i].UpdatedBy = updatedBy
	}

	if err := h.Service.AddEvents(c.Request().Context(), v); err != nil {
		return err
	}

	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}

	return c.JSON(http.StatusOK, rest.Response[[]string]{
		Message: &rest.Message{
			Text: "Events added",
		},
		Payload: ids,
	})
}

// @Summary GetEvent
// @Description GetEvent
// @Param id path string true "Event ID"
// @Success 200 {object} rest.Response[models.Event]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /events/{id} [get]
// @Tags Events
func (h *HTTP) GetEvent(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing event ID")
	}

	event, err := h.Service.GetEvent(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if event == nil {
		return echo.NewHTTPError(http.StatusNotFound, "event not found")
	}

	return c.JSON(http.StatusOK, rest.Response[models.Event]{
		Payload: *event,
	})
}

// @Summary DeleteEvent
// @Description DeleteEvent
// @Param id path string true "Event ID"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /events/{id} [delete]
// @Tags Events
func (h *HTTP) DeleteEvent(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing event ID")
	}

	if err := h.Service.RemoveEvent(c.Request().Context(), id); err != nil {
		return err
	}

	return nil
}

// @Summary DeleteEvents
// @Description DeleteEvents for multiple events
// @Param id query string true "Event ID"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /events [delete]
// @Tags Events
func (h *HTTP) DeleteEvents(c echo.Context) error {
	q, err := query.ParseWithValidator(
		c.QueryString(),
		h.Validator.DeleteEvents,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ids := q.GetValues("id")
	if len(ids) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing event ID")
	}

	if err := h.Service.RemoveEvent(c.Request().Context(), q.GetValues("id")...); err != nil {
		return err
	}

	return nil
}

// @Summary PutEvent
// @Description PutEvent
// @Param id path string true "Event ID"
// @Param body body models.Event true "Event"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /events/{id} [put]
// @Tags Events
func (h *HTTP) PutEvent(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing event ID")
	}

	v := models.Event{}
	if err := rest.BindJSON(c.Request().Body, &v); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	updatedBy := server.GetUser(c)
	v.UpdatedBy = updatedBy

	if err := h.Service.UpdateEvent(c.Request().Context(), id, &v); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rest.ResponseMessage{
		Message: &rest.Message{
			Text: "Event updated",
		},
	})
}

// /////////////////////////////////////////////////////////////
// Relations
// /////////////////////////////////////////////////////////////

// @Summary AddRelations
// @Description AddRelations
// @Param body body []models.Relation true "Relation"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /relations [post]
// @Tags Relations
func (h *HTTP) AddRelations(c echo.Context) error {
	v := []models.Relation{}
	if err := rest.BindJSONList(c.Request().Body, &v); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	updatedBy := server.GetUser(c)
	for i := range v {
		if v[i].Entity == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "missing entity")
		}

		v[i].UpdatedBy = updatedBy
	}

	if err := h.Service.AddRelations(c.Request().Context(), v); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rest.ResponseMessage{
		Message: &rest.Message{
			Text: "Relations added",
		},
	})
}

// @Summary DeleteRelations
// @Description DeleteRelations for multiple relations
// @Param entity query string true "entity"
// @Param event_id query string false "event_id"
// @Param event_group query string false "event_group"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /relations [delete]
// @Tags Relations
func (h *HTTP) DeleteRelations(c echo.Context) error {
	q, err := query.ParseWithValidator(
		c.QueryString(),
		h.Validator.DeleteRelations,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.Service.RemoveRelation(c.Request().Context(), q); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rest.ResponseMessage{
		Message: &rest.Message{
			Text: "Relation removed",
		},
	})
}

// @Summary GetRelations
// @Description GetRelations
// @Param entity query string false "entity"
// @Param event_id query string false "event_id"
// @Param event_group query string false "event_group"
// @Param sort query string false "sort"
// @Param limit query int false "limit" default(25)
// @Param offset query int false "offset"
// @Success 200 {object} rest.Response[[]models.Relation]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /relations [get]
// @Tags Relations
func (h *HTTP) GetRelations(c echo.Context) error {
	q, err := query.ParseWithValidator(c.QueryString(), h.Validator.GetRelations, query.WithDefaultLimit(DefaultLimit))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	relations, err := h.Service.GetRelations(c.Request().Context(), q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if len(relations) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "no relations found")
	}

	count, err := h.Service.GetRelationsCount(c.Request().Context(), q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, rest.Response[[]models.Relation]{
		Meta: &rest.Meta{
			TotalItemCount: count,
			Limit:          q.GetLimit(),
			Offset:         q.GetOffset(),
		},
		Payload: relations,
	})
}

// ////////////////////////////////////////////////////////////////

// @Summary Holidays
// @Description Holidays for specific date
// @Param entity query string false "entity for relation"
// @Param event_group query string false "country for relation"
// @Param date query string true "date specific event"
// @Success 200 {object} rest.Response[[]models.Event]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /holidays [get]
// @Tags Search
func (h *HTTP) Holidays(c echo.Context) error {
	q, err := query.ParseWithValidator(
		c.QueryString(),
		h.Validator.GetEventsDate,
		query.WithSkipExpressionCmp("date"),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	events, err := h.Service.GetEvents(c.Request().Context(), q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if len(events) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "no events found")
	}

	return c.JSON(http.StatusOK, rest.Response[[]models.Event]{
		Meta: &rest.Meta{
			TotalItemCount: uint64(len(events)),
			Limit:          q.GetLimit(),
			Offset:         q.GetOffset(),
		},
		Payload: events,
	})
}

// @Summary AddICS
// @Description AddICS
// @Accept multipart/form-data
// @Param file formData file true "ICS file"
// @Param event_group query string false "event_group for ics"
// @Param tz query string false "timezone like Europe/Amsterdam default UTC"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /ics [post]
// @Tags iCal
func (h *HTTP) AddICS(c echo.Context) error {
	var eventGroupNull types.Null[string]
	if eventGroup := c.QueryParam("event_group"); eventGroup != "" {
		eventGroupNull = types.NewNull(eventGroup)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to get file: "+err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open file: "+err.Error())
	}
	defer src.Close()

	tz := strings.TrimSpace(c.QueryParam("tz"))
	defaultTZ := time.UTC
	if tz != "" {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid timezone: "+tz+" "+err.Error())
		}

		defaultTZ = loc
	}

	if err := h.Service.AddIcal(c.Request().Context(), src, defaultTZ, eventGroupNull, server.GetUser(c)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to add ICS: "+err.Error())
	}

	return c.JSON(http.StatusOK, rest.ResponseMessage{
		Message: &rest.Message{
			Text: "ICS added",
		},
	})
}

// @Summary GetICS
// @Description GetICS
// @Param entity query string false "entity for relation"
// @Param event_group query string false "country"
// @Param year query string false "specific year events"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /ics [get]
// @Tags iCal
func (h *HTTP) GetICS(c echo.Context) error {
	q, err := query.ParseWithValidator(
		c.QueryString(),
		h.Validator.GetICS,
		query.WithSkipExpressionCmp("year"),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	events, err := h.Service.GetEventsICS(c.Request().Context(), q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// convert ics format
	category := strings.Join(q.GetValues("entity"), ",")
	fileName := strings.ToLower(strings.ReplaceAll(category, ",", "_"))
	if fileName == "" {
		fileName = "events"
	}
	if category == "" {
		category = "Holidays"
	}

	str, err := ical.GenerateICS(events, category)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// send ics file
	c.Response().Header().Set(echo.HeaderContentType, "text/calendar")
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+fileName+".ics")
	c.Response().WriteHeader(http.StatusOK)

	_, err = c.Response().Write([]byte(str))

	return err
}
