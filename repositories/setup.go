package repositories

import "database/sql"

type Setup interface {
	CreateMigrationTable(tableName string) (sql.Result, error)
}

var _ Setup = &postgresSetup{}
