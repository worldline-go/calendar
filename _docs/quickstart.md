# Getting Started

This service is written in [Go](https://golang.org/) and does not require any additional dependencies to run.  
Storage need [PostgreSQL](https://www.postgresql.org/).

## Install

### Binary

Download the latest release from [GitHub](https://github.com/worldline-go/calendar/releases/latest)

Extract it from archive and before to run, you need to have configuration file.

## Configuration

Give the path to the configuration file using `CONFIG_FILE` environment variable or use default file name `calendar.[toml|yaml|yml|json]` in the current directory.

```yaml
log_level: info
port: 8080

db_type: pgx
db_datasource: postgres://postgres@localhost:5432/postgres?sslmode=disable # default is empty
db_schema: public

migrate:
  db_datasource: postgres://postgres@localhost:5432/postgres?sslmode=disable # default is empty
  db_type: pgx
  db_schema: public
  db_table: calendar_migrations
```

> Configuration migration's connect and database's connect are separated.
