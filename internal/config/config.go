package config

import (
	"context"
	"fmt"

	"github.com/worldline-go/igconfig"
	"github.com/worldline-go/igconfig/loader"
	"github.com/worldline-go/logz"
	"github.com/worldline-go/tell"
)

var (
	ServiceName    = "calendar"
	ServiceVersion = "v0.0.0"
	ServiceDomain  = "admin"

	ServiceLog = ServiceName + "@" + ServiceVersion
)

// Config contains application configuration for this command.
type Config struct {
	LogLevel string `cfg:"log_level" default:"info"`
	Port     uint   `cfg:"port"      default:"8080"`

	DBType       string `cfg:"db_type"       default:"pgx"`
	DBDataSource string `cfg:"db_datasource" log:"false"`
	DBSchema     string `cfg:"db_schema"     default:"public"`

	Migrate Migrate `cfg:"migrate"`

	Telemetry tell.Config
}

// Migrate contains database connection to run the migrations.
type Migrate struct {
	DBDatasource string `cfg:"db_datasource" log:"false"`
	DBType       string `cfg:"db_type"       default:"pgx"`
	DBSchema     string `cfg:"db_schema"     default:"public"`
	DBTable      string `cfg:"db_table"      default:"calendar_migrations"`
}

func init() {
	loader.VaultSecretAdditionalPaths = append(loader.VaultSecretAdditionalPaths,
		loader.AdditionalPath{Map: "migrate", Name: "migrations"},
	)
}

func Load(ctx context.Context) (*Config, error) {
	cfg := &Config{}

	if err := igconfig.LoadConfigWithContext(ctx, ServiceDomain+"/"+ServiceName, cfg); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	if err := logz.SetLogLevel(cfg.LogLevel); err != nil {
		return nil, fmt.Errorf("parse log level %s: %w", cfg.LogLevel, err)
	}

	return cfg, nil
}
