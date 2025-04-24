package main

import (
	"context"
	"fmt"

	"github.com/rakunlabs/chu"
	"github.com/rs/zerolog/log"
	"github.com/worldline-go/initializer"
	"github.com/worldline-go/tell"

	"github.com/worldline-go/calendar/internal/adapter/repository"
	"github.com/worldline-go/calendar/internal/config"
	"github.com/worldline-go/calendar/internal/core/service"
	"github.com/worldline-go/calendar/internal/server"
)

var (
	version = "v0.0.0"
	commit  = ""
	date    = ""
)

func main() {
	config.ServiceVersion = version

	initializer.Init(
		run,
		initializer.WithMsgf("%s [%s] build %s %s", config.ServiceName, config.ServiceVersion, commit, date),
	)
}

func run(ctx context.Context) error {
	cfg, err := config.Load(ctx)
	if err != nil {
		return err
	}

	log.Log().RawJSON("config", chu.MarshalJSON(cfg)).Msg("loaded Config")

	// ///////////////////////////////////////////////////////
	// telemetry initialize
	collector, err := tell.New(ctx, cfg.Telemetry)
	if err != nil {
		return fmt.Errorf("failed to init telemetry; %w", err)
	}
	// flush metrics on failure
	defer collector.Shutdown()

	// ///////////////////////////////////////////////////////
	// database operations
	if err := repository.MigrateDB(ctx, cfg); err != nil {
		return fmt.Errorf("failed database migration: %w", err)
	}

	calendarPostgresAdapter, err := repository.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// ///////////////////////////////////////////////////////
	// service initialize
	svc, err := service.NewCalendarService(ctx, calendarPostgresAdapter)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	// ///////////////////////////////////////////////////////
	// server initialize
	srv, err := server.NewServer(ctx, svc)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// ///////////////////////////////////////////////////////
	// start server
	initializer.ShutdownAdd(srv.Stop, "server")

	return srv.Start(fmt.Sprintf(":%d", cfg.Port))
}
