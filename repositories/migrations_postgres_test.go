package repositories_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	_ "github.com/lib/pq"
	"github.com/raytung/g6/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/raytung/g6/pkg/tests/docker"
	"strings"
)

func Test_Integration_Migrations_Postgres_CreateTable(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	out, err, tearDown := docker.Cli(&docker.Options{
		Command:       "run",
		ContainerName: "g6_itest_create_migration_table",
		Image:         "postgres:alpine",
		Publish:       "5435:5432",
		Env: map[string]string{
			"POSTGRES_USER":     "g6_test",
			"POSTGRES_DB":       "g6_test",
			"POSTGRES_PASSWORD": "password",
		},
	})
	defer tearDown()
	assert.NoError(t, err, string(out))

	db, err := sql.Open("postgres", "postgres://g6_test:password@0.0.0.0:5435/g6_test?sslmode=disable")
	assert.NoError(t, err)

	docker.WaitForDB(t, db)

	pg := repositories.NewPostgresMigrations(db, "g6_migrations")
	type args struct {
		tableName string
	}
	tests := []struct {
		name          string
		args          args
		expectedError error
	}{
		{
			name:          "creates migration table",
			args:          args{"g6_migrations"},
			expectedError: nil,
		},

		{
			name:          "does not create migration table if already exist",
			args:          args{"g6_migrations"},
			expectedError: errors.New("pq: relation \"g6_migrations\" already exists"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pg.CreateTable(tt.args.tableName)
			assert.Equal(t, fmt.Sprintf("%v", tt.expectedError), fmt.Sprintf("%v", err))
		})
	}
}

func Test_Integration_Migrations_Postgres_TableExists(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	out, err, tearDown := docker.Cli(&docker.Options{
		Command:       "run",
		ContainerName: "g6_itest_migrations_table_exists",
		Image:         "postgres:alpine",
		Publish:       "5436:5432",
		Env: map[string]string{
			"POSTGRES_USER":     "g6_test",
			"POSTGRES_DB":       "g6_test",
			"POSTGRES_PASSWORD": "password",
		},
	})

	defer tearDown()
	assert.NoError(t, err, string(out))

	db, err := sql.Open("postgres", "postgres://g6_test:password@0.0.0.0:5436/g6_test?sslmode=disable")
	assert.NoError(t, err)

	docker.WaitForDB(t, db)

	pg := repositories.NewPostgresMigrations(db, "g6_migrations")

	_, err = pg.CreateTable("g6_migrations")
	assert.NoError(t, err)

	tests := []struct {
		name          string
		table         string
		exists        bool
		expectedError error
	}{
		{
			name:          "table exists",
			table:         "g6_migrations",
			exists:        true,
			expectedError: nil,
		},

		{
			name:          "table does not exist",
			table:         "random_table",
			exists:        false,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pg.TableExists(tt.table)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.exists, got)
		})
	}
}

func Test_Integration_Migrations_Postgres_Run(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	out, err, tearDown := docker.Cli(&docker.Options{
		Command:       "run",
		ContainerName: "g6_itest_migrations_run",
		Image:         "postgres:alpine",
		Publish:       "5437:5432",
		Env: map[string]string{
			"POSTGRES_USER":     "g6_test",
			"POSTGRES_DB":       "g6_test",
			"POSTGRES_PASSWORD": "password",
		},
	})

	defer tearDown()
	assert.NoError(t, err, string(out))

	db, err := sql.Open("postgres", "postgres://g6_test:password@0.0.0.0:5437/g6_test?sslmode=disable")
	assert.NoError(t, err)

	docker.WaitForDB(t, db)

	pg := repositories.NewPostgresMigrations(db, "g6_migrations")

	_, err = pg.CreateTable("g6_migrations")
	assert.NoError(t, err)

	tests := []struct {
		name          string
		migration     *repositories.Migration
		expectedError error
		expectedTable string
	}{
		{
			name:          "run migrations",
			expectedError: nil,
			expectedTable: "users",
			migration: &repositories.Migration{
				Name: "V1234__create_users_table",
				Query: strings.Join([]string{
					"CREATE TABLE users (",
					"	id SERIAL PRIMARY KEY,",
					"	email VARCHAR UNIQUE NOT NULL",
					");",
				}, "\n"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pg.Run(tt.migration)
			assert.Equal(t, tt.expectedError, err)
			rows, err := db.Query("SELECT COUNT(*) FROM " + tt.expectedTable)
			assert.NoError(t, err)
			assert.True(t, rows.Next())
			var count int
			assert.NoError(t, rows.Scan(&count))
			assert.Equal(t, 0, count)
			migrationsTableRows, err := db.Query("SELECT COUNT(*) FROM g6_migrations WHERE name = $1", tt.migration.Name)
			var migrationsTableCount int
			assert.NoError(t, err)
			assert.True(t, migrationsTableRows.Next())
			assert.NoError(t, migrationsTableRows.Scan(&migrationsTableCount))
			assert.Equal(t, 1, migrationsTableCount)
		})
	}
}
