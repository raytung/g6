package repositories

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"time"
)

var _ Migrations = &pgMigrations{}

func createMigrationQuery(table string) string {
	return fmt.Sprintf(`
CREATE TABLE "%s" (
	id          SERIAL    PRIMARY KEY,
    Name        VARCHAR   NOT NULL UNIQUE,
    migrated_at TIMESTAMP NOT NULL
);`, table);
}

type pgMigrations struct {
	db    *sql.DB
	table string
}

func (pg *pgMigrations) CreateTable(tableName string) (sql.Result, error) {
	return pg.db.Exec(createMigrationQuery(tableName))
}

func (pg *pgMigrations) TableExists(table string) (bool, error) {
	query := fmt.Sprintf(`
SELECT *
FROM "%s"
LIMIT 1
`, table)
	_, err := pg.db.Query(query)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "undefined_table" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (pg *pgMigrations) Run(migration *Migration) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(migration.Query)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(fmt.Sprintf("INSERT INTO \"%s\" (name, migrated_at) VALUES ($1, $2)", pg.table), migration.Name, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func NewPostgresMigrations(conn *sql.DB, table string) *pgMigrations {
	return &pgMigrations{conn, table}
}
