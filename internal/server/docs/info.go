package docs

import (
	_ "github.com/swaggo/swag"
	"github.com/worldline-go/swagger"

	"github.com/worldline-go/calendar/internal/config"
)

func Info() error {
	return swagger.SetInfo(
		swagger.WithTitle(config.ServiceName),
		swagger.WithVersion(config.ServiceVersion),
	)
}
