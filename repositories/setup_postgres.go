package repositories

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
)

var _ Setup = &postgresSetup{}

func createMigrationQuery(table string) string {
	return fmt.Sprintf(`
CREATE TABLE "%s" (
	id          SERIAL    PRIMARY KEY,
    name        VARCHAR   NOT NULL UNIQUE,
    migrated_at TIMESTAMP NOT NULL
);`, table);
}

type postgresSetup struct {
	db *sql.DB
}

func (pg *postgresSetup) CreateMigrationTable(tableName string) (sql.Result, error) {
	return pg.db.Exec(createMigrationQuery(tableName))
}

func NewPostgresSetup() *postgresSetup {
	return &postgresSetup{}
}
