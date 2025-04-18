package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
	GetEvents     *query.Validator
	GetEventsDate *query.Validator
	GetICS        *query.Validator
}

func NewHTTP(svc *service.Service) (*HTTP, error) {
	validatorGetEvents, err := query.NewValidator(
		query.WithField(query.WithNotAllowed()),
		query.WithValues(query.WithIn("id", "code", "country", "name", "description", "disabled")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator for GetEvents: %w", err)
	}

	validatorGetEventsDate, err := query.NewValidator(
		query.WithField(query.WithNotAllowed()),
		query.WithValues(query.WithIn("code", "country", "date")),
		query.WithValue("code", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithValue("country", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithValue("date", query.WithOperator(query.OperatorEq), query.WithNotEmpty()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator for GetEventsDate: %w", err)
	}

	validatorGetICS, err := query.NewValidator(
		query.WithField(query.WithNotAllowed()),
		query.WithValues(query.WithIn("code", "country", "year", "tz")),
		query.WithValue("tz", query.WithOperator(query.OperatorEq)),
		query.WithValue("code", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithValue("country", query.WithOperator(query.OperatorEq, query.OperatorIn)),
		query.WithValue("year", query.WithOperator(query.OperatorEq, query.OperatorIn)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator for GetICS: %w", err)
	}

	return &HTTP{
		Service: svc,
		Validator: QueryValidator{
			GetEvents:     validatorGetEvents,
			GetEventsDate: validatorGetEventsDate,
			GetICS:        validatorGetICS,
		},
	}, nil
}

func (h *HTTP) RegisterRoutes(g *echo.Group) {
	g.GET("/events", h.GetEvents)
	g.POST("/events", h.AddEvents)

	g.GET("/events/{id}", h.GetEvent)
	g.DELETE("/events/{id}", h.RemoveEvent)
	g.PUT("/events/{id}", h.RemoveEvent)
	g.PATCH("/events/{id}", h.RemoveEvent)

	g.GET("/relations", h.GetRelations)
	g.POST("/relations", h.AddRelations)

	g.GET("/relations/{id}", h.GetRelation)
	g.DELETE("/relations/{id}", h.RemoveRelation)

	g.GET("/workday", h.WorkDay)
	g.POST("/ics", h.AddICS)
	g.GET("/ics", h.GetICS)
}

// @Summary GetEvents
// @Description GetEvents
// @Param id query string false "id"
// @Param name query string false "name"
// @Param description query string false "description"
// @Param disabled query bool false "disabled"
// @Param code query int false "code for relation"
// @Param country query string false "country for relation"
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
		query.WithDefaultLimit(25),
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
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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

// @Summary RemoveEvent
// @Description RemoveEvent
// @Param id path string true "Event ID"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /events/{id} [delete]
// @Tags Events
func (h *HTTP) RemoveEvent(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing event ID")
	}

	if err := h.Service.RemoveEvent(c.Request().Context(), id); err != nil {
		return err
	}

	return nil
}

// /////////////////////////////////////////////////////////////
// Relations
// /////////////////////////////////////////////////////////////

// @Summary AddRelations
// @Description AddRelations
// @Param body body []models.Relation true "Relation"
// @Success 200 {object} rest.Response[[]string]
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
		v[i].UpdatedBy = updatedBy
	}

	if err := h.Service.AddRelations(c.Request().Context(), v); err != nil {
		return err
	}

	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}

	return c.JSON(http.StatusOK, rest.Response[[]string]{
		Message: &rest.Message{
			Text: "Relations added",
		},
		Payload: ids,
	})
}

// @Summary RemoveRelation
// @Description RemoveRelation
// @Param id path string true "Relation ID"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /relations/{id} [delete]
// @Tags Relations
func (h *HTTP) RemoveRelation(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing relation ID")
	}

	if err := h.Service.RemoveRelation(c.Request().Context(), id); err != nil {
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
// @Param id query string false "id"
// @Param code query int false "code for relation"
// @Param country query string false "country for relation"
// @Param limit query int false "limit" default(25)
// @Param offset query int false "offset"
// @Success 200 {object} rest.Response[[]models.Relation]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /relations [get]
// @Tags Relations
func (h *HTTP) GetRelations(c echo.Context) error {
	q, err := query.ParseWithValidator(c.QueryString(), h.Validator.GetEvents)
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
			TotalItemCount: uint64(count),
			Limit:          q.GetLimit(),
			Offset:         q.GetOffset(),
		},
		Payload: relations,
	})
}

// @Summary GetRelation
// @Description GetRelation
// @Param id path string true "Relation ID"
// @Success 200 {object} rest.Response[models.Relation]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /relations/{id} [get]
// @Tags Relations
func (h *HTTP) GetRelation(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing relation ID")
	}

	relation, err := h.Service.GetRelation(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if relation == nil {
		return echo.NewHTTPError(http.StatusNotFound, "relation not found")
	}

	return c.JSON(http.StatusOK, rest.Response[models.Relation]{
		Payload: *relation,
	})
}

// ////////////////////////////////////////////////////////////////

// @Summary WorkDay
// @Description GetEvents for specific date
// @Param code query int false "code for relation"
// @Param country query string false "country for relation"
// @Param date query string true "date specific event"
// @Success 200 {object} rest.Response[[]models.Event]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /workday [get]
// @Tags Search
func (h *HTTP) WorkDay(c echo.Context) error {
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
// @Param code query int false "code for relation"
// @Param country query string false "country for relation"
// @Param tz query string false "timezone like Europe/Amsterdam"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /ics [post]
// @Tags iCal
func (h *HTTP) AddICS(c echo.Context) error {
	var codeNull types.Null[int64]
	if code := c.QueryParam("code"); code != "" {
		codeInt, err := strconv.ParseInt(code, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid code: "+err.Error())
		}

		codeNull = types.NewNull(codeInt)
	}

	var countryNull types.Null[string]
	if country := c.QueryParam("country"); country != "" {
		countryNull = types.NewNull(country)
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

	commonRelation := models.Relation{
		Code:      codeNull,
		Country:   countryNull,
		UpdatedBy: server.GetUser(c),
	}

	tz := strings.TrimSpace(c.QueryParam("tz"))
	defaultTZ := time.UTC
	if tz != "" {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid timezone: "+tz+" "+err.Error())
		}

		defaultTZ = loc
	}

	if err := h.Service.AddIcal(c.Request().Context(), src, commonRelation, defaultTZ); err != nil {
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
// @Param code query int false "code for relation"
// @Param country query string false "country for relation"
// @Param year query int true "specific year events"
// @Param tz query string false "timezone like Europe/Amsterdam"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /ics [get]
// @Tags iCal
func (h *HTTP) GetICS(c echo.Context) error {
	q, err := query.ParseWithValidator(
		c.QueryString(),
		h.Validator.GetICS,
		query.WithSkipExpressionCmp("year", "tz"),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var tzLoc *time.Location
	if tzCmp, _ := q.Values["tz"]; len(tzCmp) > 0 {
		tz := tzCmp[0].Value.(string)
		tz, _ = url.QueryUnescape(tz)
		loc, err := time.LoadLocation(tz)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid timezone: "+tz+" "+err.Error())
		}

		tzLoc = loc
	}

	events, err := h.Service.GetEvents(c.Request().Context(), q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if tzLoc != nil {
		for i := range events {
			events[i].DateFrom = types.Time{Time: events[i].DateFrom.In(tzLoc)}
			events[i].DateTo = types.Time{Time: events[i].DateTo.In(tzLoc)}
		}
	}

	// convert ics format
	str, err := ical.GenerateICS(events)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// send ics file
	c.Response().Header().Set(echo.HeaderContentType, "text/calendar")
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=events.ics")
	c.Response().WriteHeader(http.StatusOK)

	_, err = c.Response().Write([]byte(str))

	return err
}
