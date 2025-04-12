package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/worldline-go/query"
	"github.com/worldline-go/rest"
	"github.com/worldline-go/rest/server"

	"github.com/worldline-go/calendar/internal/service"
	"github.com/worldline-go/calendar/pkg/models"
)

type HTTP struct {
	Service *service.Service

	Validator QueryValidator
}

type QueryValidator struct {
	GetHolidays *query.Validator
}

func NewHTTP(svc *service.Service) (*HTTP, error) {
	validatorGetHolidays, err := query.NewValidator(
		query.WithField(query.WithNotAllowed()),
		query.WithValues(query.WithIn("id", "code", "country", "name", "description", "disabled", "date", "year")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator for GetHolidays: %w", err)
	}

	return &HTTP{
		Service: svc,
		Validator: QueryValidator{
			GetHolidays: validatorGetHolidays,
		},
	}, nil
}

func (h *HTTP) RegisterRoutes(g *echo.Group) {
	g.GET("/holidays", h.GetHolidays)
	g.POST("/holidays", h.AddHoliday)
	g.POST("/holidays-bulk", h.AddHolidayBulk)

	g.GET("/holidays/{id}", h.GetHoliday)
	g.DELETE("/holidays/{id}", h.RemoveHoliday)
	g.PUT("/holidays/{id}", h.RemoveHoliday)
	g.PATCH("/holidays/{id}", h.RemoveHoliday)

	g.GET("/relations", h.GetRelations)
	g.POST("/relations", h.AddRelation)
	g.POST("/relations-bulk", h.AddRelationBulk)

	g.GET("/relations/{id}", h.GetRelation)
	g.DELETE("/relations/{id}", h.RemoveRelation)

	g.GET("/workday", h.WorkDay)
}

// @Summary GetHolidays
// @Description GetHolidays
// @Param id query string false "id"
// @Param name query string false "name"
// @Param description query string false "description"
// @Param disabled query bool false "disabled"
// @Param code query int false "code for relation"
// @Param country query string false "country for relation"
// @Param date query string false "date specific holiday"
// @Param year query int false "year"
// @Param limit query int false "limit" default(25)
// @Param offset query int false "offset"
// @Success 200 {object} rest.Response[[]models.Holiday]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /holidays [get]
func (h *HTTP) GetHolidays(c echo.Context) error {
	q, err := query.ParseWithValidator(
		c.QueryString(),
		h.Validator.GetHolidays,
		query.WithSkipExpressionCmp("date", "year"),
		query.WithDefaultLimit(25),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	holidays, err := h.Service.GetHolidays(c.Request().Context(), q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if len(holidays) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "no holidays found")
	}

	count, err := h.Service.GetHolidaysCount(c.Request().Context(), q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, rest.Response[[]*models.Holiday]{
		Meta: &rest.Meta{
			TotalItemCount: uint64(count),
			Limit:          q.GetLimit(),
			Offset:         q.GetOffset(),
		},
		Payload: holidays,
	})
}

// @Summary AddHoliday
// @Description AddHoliday
// @Param body body models.Holiday true "Holiday"
// @Success 200 {object} rest.Response[string]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /holidays [post]
func (h *HTTP) AddHoliday(c echo.Context) error {
	v := models.Holiday{}
	if err := c.Bind(&v); err != nil {
		return err
	}

	v.UpdatedBy = server.GetUser(c)

	if err := h.Service.AddHoliday(c.Request().Context(), &v); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rest.Response[string]{
		Message: &rest.Message{
			Text: "Holiday added",
		},
		Payload: v.ID,
	})
}

// @Summary AddHolidayBulk
// @Description AddHolidayBulk
// @Param body body []models.Holiday true "Holiday"
// @Success 200 {object} rest.Response[[]string]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /holidays-bulk [post]
func (h *HTTP) AddHolidayBulk(c echo.Context) error {
	v := []*models.Holiday{}
	if err := c.Bind(&v); err != nil {
		return err
	}

	updatedBy := server.GetUser(c)
	for i := range v {
		v[i].UpdatedBy = updatedBy
	}

	if err := h.Service.AddHoliday(c.Request().Context(), v...); err != nil {
		return err
	}

	ids := make([]string, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}

	return c.JSON(http.StatusOK, rest.Response[[]string]{
		Message: &rest.Message{
			Text: "Holidays added",
		},
		Payload: ids,
	})
}

// @Summary GetHoliday
// @Description GetHoliday
// @Param id path string true "Holiday ID"
// @Success 200 {object} rest.Response[models.Holiday]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /holidays/{id} [get]
func (h *HTTP) GetHoliday(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing holiday ID")
	}

	holiday, err := h.Service.GetHoliday(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if holiday == nil {
		return echo.NewHTTPError(http.StatusNotFound, "holiday not found")
	}

	return c.JSON(http.StatusOK, rest.Response[models.Holiday]{
		Payload: *holiday,
	})
}

// @Summary RemoveHoliday
// @Description RemoveHoliday
// @Param id path string true "Holiday ID"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /holidays/{id} [delete]
func (h *HTTP) RemoveHoliday(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing holiday ID")
	}

	if err := h.Service.RemoveHoliday(c.Request().Context(), id); err != nil {
		return err
	}

	return nil
}

// /////////////////////////////////////////////////////////////
// Relations
// /////////////////////////////////////////////////////////////

// @Summary AddRelation
// @Description AddRelation
// @Param body body models.Relation true "Relation"
// @Success 200 {object} rest.Response[string]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /relations [post]
func (h *HTTP) AddRelation(c echo.Context) error {
	v := models.Relation{}
	if err := c.Bind(&v); err != nil {
		return err
	}

	v.UpdatedBy = server.GetUser(c)

	if err := h.Service.AddRelation(c.Request().Context(), &v); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rest.Response[string]{
		Message: &rest.Message{
			Text: "Relation added",
		},
		Payload: v.ID,
	})
}

// @Summary AddRelationBulk
// @Description AddRelationBulk
// @Param body body []models.Relation true "Relation"
// @Success 200 {object} rest.Response[[]string]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /relations-bulk [post]
func (h *HTTP) AddRelationBulk(c echo.Context) error {
	v := []*models.Relation{}
	if err := c.Bind(&v); err != nil {
		return err
	}

	updatedBy := server.GetUser(c)
	for i := range v {
		v[i].UpdatedBy = updatedBy
	}

	if err := h.Service.AddRelation(c.Request().Context(), v...); err != nil {
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
// @Param id path string true "Holiday ID"
// @Success 200 {object} rest.ResponseMessage
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /relations/{id} [delete]
func (h *HTTP) RemoveRelation(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing holiday ID")
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
func (h *HTTP) GetRelations(c echo.Context) error {
	q, err := query.ParseWithValidator(c.QueryString(), h.Validator.GetHolidays)
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

	return c.JSON(http.StatusOK, rest.Response[[]*models.Relation]{
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
// @Param id path string true "Holiday ID"
// @Success 200 {object} rest.Response[models.Relation]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /relations/{id} [get]
func (h *HTTP) GetRelation(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing holiday ID")
	}

	relation, err := h.Service.GetRelation(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if relation == nil {
		return echo.NewHTTPError(http.StatusNotFound, "holiday not found")
	}

	return c.JSON(http.StatusOK, rest.Response[models.Relation]{
		Payload: *relation,
	})
}

// ////////////////////////////////////////////////////////////////

// @Summary WorkDay
// @Description WorkDay
// @Param date query string true "date"
// @Param country query string false "country"
// @Param code query int false "code for relation"
// @Success 200 {object} rest.Response[string]
// @Failure 400 {object} rest.ResponseMessage
// @Failure 500 {object} rest.ResponseMessage
// @Router /workday [get]
func (h *HTTP) WorkDay(c echo.Context) error {

	return nil
}
