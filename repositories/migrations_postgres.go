package repositories

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
)

var _ Migrations = &pgMigrations{}

func createMigrationQuery(table string) string {
	return fmt.Sprintf(`
CREATE TABLE "%s" (
	id          SERIAL    PRIMARY KEY,
    name        VARCHAR   NOT NULL UNIQUE,
    migrated_at TIMESTAMP NOT NULL
);`, table);
}

type pgMigrations struct {
	db *sql.DB
}

func (pg *pgMigrations) CreateTable(tableName string) (sql.Result, error) {
	return pg.db.Exec(createMigrationQuery(tableName))
}

func NewPostgresMigrations(conn *sql.DB) *pgMigrations {
	return &pgMigrations{conn}
}
