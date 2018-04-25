package repositories_test

import (
	"database/sql"
	"errors"
	"fmt"
	"os/exec"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/raytung/g6/repositories"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_Migrations_Postgres_CreateTable(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	container := "g6_itest_create_migration_table"

	dbPort := "5435"
	out, err, tearDown := startPostgres(container, dbPort)
	defer tearDown()
	assert.NoError(t, err, string(out))

	db, err := sql.Open("postgres", "postgres://g6_test:password@0.0.0.0:5435/g6_test?sslmode=disable")
	assert.NoError(t, err)

	waitForPostgres(t, db)

	pg := repositories.NewPostgresMigrations(db)
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

func startPostgres(container string, port string) ([]byte, error, func() ([]byte, error)) {
	out, err := exec.Command("docker", "run",
		"--rm",
		"--name", container,
		"--detach",
		"--publish", port+":5432",
		"--env", "POSTGRES_USER=g6_test",
		"--env", "POSTGRES_DB=g6_test",
		"--env", "POSTGRES_PASSWORD=password",
		"postgres:alpine",
	).Output()

	tearDown := func() ([]byte, error) {
		return exec.Command("docker", "stop", container).Output()
	}

	return out, err, tearDown
}

func waitForPostgres(t *testing.T, db *sql.DB) {
	t.Helper()
	timeout := 30 * time.Second
	elapsed := 0 * time.Second
	waitTime := 500 * time.Millisecond

	for {
		time.Sleep(waitTime)
		elapsed += waitTime

		if elapsed >= timeout {
			t.Fatalf("Cannot connect to postgres after %v", timeout)
		}

		if err := db.Ping(); err != nil {
			if testing.Verbose() {
				fmt.Printf("DB Ping error: %+v\n", err)
			}
			continue
		}

		break
	}
}

func Test_Integration_Migrations_Postgres_TableExists(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	container := "g6_itest_migrations_table_exists"

	dbPort := "5436"
	out, err, tearDown := startPostgres(container, dbPort)
	defer tearDown()
	assert.NoError(t, err, string(out))

	db, err := sql.Open("postgres", "postgres://g6_test:password@0.0.0.0:5436/g6_test?sslmode=disable")
	assert.NoError(t, err)

	waitForPostgres(t, db)

	pg := repositories.NewPostgresMigrations(db)

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
