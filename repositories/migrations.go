package repositories

import "database/sql"

type Migration struct {
	Name  string
	Query string
}

type Options struct {
	table string
}

type MigrationsRunner interface {
	Run(migration *Migration) error
}

type Migrations interface {
	CreateTable(table string) (sql.Result, error)
	TableExists(table string) (bool, error)
}

var _ Migrations = &pgMigrations{}
