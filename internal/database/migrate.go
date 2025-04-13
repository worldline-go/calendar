package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/worldline-go/igmigrator/v2"

	"github.com/worldline-go/calendar/internal/config"
)

//go:embed migrations/*
var migrationFS embed.FS

func MigrateDB(ctx context.Context, cfg *config.Config) error {
	if cfg.Migrate.DBDatasource == "" {
		return fmt.Errorf("migrate database datasource is empty")
	}

	migration, err := fs.Sub(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("migrate database fs sub: %w", err)
	}

	db, err := sqlx.Connect(cfg.Migrate.DBType, cfg.Migrate.DBDatasource)
	if err != nil {
		return fmt.Errorf("migrate database connect: %w", err)
	}

	defer db.Close()

	result, err := igmigrator.Migrate(ctx, db, &igmigrator.Config{
		Migrations:     migration,
		Schema:         cfg.Migrate.DBSchema,
		MigrationTable: cfg.Migrate.DBTable,
	})
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	for mPath, m := range result.Path {
		if m.NewVersion != m.PrevVersion {
			log.Info().Msgf("ran migration [%s] from version [%d] to [%d]", mPath, m.PrevVersion, m.NewVersion)
		}
	}

	return nil
}
