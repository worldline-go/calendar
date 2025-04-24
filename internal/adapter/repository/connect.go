package repository

import (
	"context"
	"fmt"
	"time"

	// Register pgx driver for SQL.
	"github.com/doug-martin/goqu/v9"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/worldline-go/calendar/internal/config"
)

var (
	ConnMaxLifetime = 15 * time.Minute
	MaxIdleConns    = 3
	MaxOpenConns    = 3
)

type Database struct {
	q *goqu.Database
}

// New attempts to connect to database server and returns a new Database instance.
func New(ctx context.Context, cfg *config.Config) (*Database, error) {
	db, err := sqlx.ConnectContext(ctx, cfg.DBType, cfg.DBDataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetConnMaxLifetime(ConnMaxLifetime)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)

	return newDB(db, cfg.DBSchema), nil
}

func newDB(db *sqlx.DB, schema string) *Database {
	setSchema(schema)

	return &Database{
		q: goqu.New("postgres", db),
	}
}
