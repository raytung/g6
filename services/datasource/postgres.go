package datasource

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type DBQuery interface {
	Query(string, ...interface{}) (DBRows, error)
}

type DBExec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type DBQueryExec interface {
	DBQuery
	DBExec
}

type DBRows interface {
	Next() bool
}

type DBRowsScanner interface {
	Scan(...interface{}) error
}

var _ DBQuery = &postgres{}
var _ DBExec = &postgres{}

type postgres struct {
	db *sql.DB
}

func (pg *postgres) Query(query string, args ...interface{}) (DBRows, error) {
	return pg.db.Query(query, args...)
}

func (pg *postgres) Exec(query string, args ...interface{}) (sql.Result, error) {
	return pg.db.Exec(query, args...)
}

func NewPostgres(db *sql.DB) *postgres {
	return &postgres{
		db,
	}
}
