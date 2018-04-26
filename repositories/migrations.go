package repositories

import (
	"database/sql"
	"time"
)

type Migration struct {
	Name  string
	Query string
}

type MigrationQueryResult struct {
	HasResults bool
	ID         int
	Name       string
	MigratedAt time.Time
}

type Options struct {
	table string
}

type MigrationsRunner interface {
	Run(migration *Migration) error
}

type MigrationsLatestInfo interface {
	Latest() (*Migration, error)
}

type Migrations interface {
	CreateTable(table string) (sql.Result, error)
	TableExists(table string) (bool, error)
	Latest() (*MigrationQueryResult, error)
}

var _ Migrations = &pgMigrations{}
