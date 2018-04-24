package repositories

import "database/sql"

type Migrations interface {
	CreateTable(tableName string) (sql.Result, error)
}

var _ Migrations = &pgMigrations{}
