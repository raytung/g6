package repositories

import "database/sql"

type Migrations interface {
	CreateTable(table string) (sql.Result, error)
	TableExists(table string) (bool, error)
}

var _ Migrations = &pgMigrations{}
