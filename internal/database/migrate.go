package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/worldline-go/igmigrator/v2"

	"github.com/worldline-go/calendar/internal/config"
)

func MigrateDB(ctx context.Context, cfg *config.Config) error {
	if cfg.Migrate.DBDatasource == "" {
		return fmt.Errorf("migrate database datasource is empty")
	}

	db, err := sqlx.Connect(cfg.Migrate.DBType, cfg.Migrate.DBDatasource)
	if err != nil {
		return fmt.Errorf("migrate database connect: %w", err)
	}

	defer db.Close()

	result, err := igmigrator.Migrate(ctx, db, &igmigrator.Config{
		MigrationsDir:  "migrations",
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
