package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/igconfig"
	"github.com/worldline-go/initializer"
	"github.com/worldline-go/tell"

	"github.com/worldline-go/calendar/internal/config"
	"github.com/worldline-go/calendar/internal/database"
	"github.com/worldline-go/calendar/internal/server"
	"github.com/worldline-go/calendar/internal/service"
)

func main() {
	initializer.Init(
		run,
		initializer.WithMsgf("%s [%s]", config.ServiceName, config.ServiceVersion),
	)
}

func run(ctx context.Context) error {
	cfg, err := config.Load(ctx)
	if err != nil {
		return err
	}

	log.Log().Object("config", igconfig.Printer{Value: cfg}).Msg("loaded Config")

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
	if err := database.MigrateDB(ctx, cfg); err != nil {
		return fmt.Errorf("failed database migration: %w", err)
	}

	db, err := database.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// ///////////////////////////////////////////////////////
	// service initialize
	svc, err := service.New(ctx, db)
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
