package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/worldline-go/echo-swagger"
	"github.com/worldline-go/rest/server"

	"github.com/worldline-go/calendar/internal/adapter/handler"
	"github.com/worldline-go/calendar/internal/config"
	"github.com/worldline-go/calendar/internal/core/port"
	"github.com/worldline-go/calendar/internal/server/docs"
)

// @title calendar API
// @BasePath /calendar/v1
func NewServer(ctx context.Context, svc port.CalendarService) (*server.Server, error) {
	return server.New(ctx, config.ServiceName, func(ctx context.Context, mux *http.ServeMux, e *echo.Echo) error {
		if err := docs.Info(); err != nil {
			log.Warn().Err(err).Msg("failed to set swagger info")
		}

		handleHTTP, err := handler.NewHTTP(svc)
		if err != nil {
			return fmt.Errorf("failed to create http handler: %w", err)
		}

		sGroup := e.Group("/calendar")
		sGroup.Any("/swagger/*", echoSwagger.WrapHandler)

		v1Group := sGroup.Group("/v1")

		// ////////////////////////////

		handleHTTP.RegisterRoutes(v1Group)

		// ////////////////////////////
		// add handler to mux
		mux.HandleFunc("/calendar/", e.ServeHTTP)

		return nil
	})
}
