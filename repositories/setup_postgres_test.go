package repositories

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"database/sql"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func Test_Integration_Setup_Postgres_CreateMigrationTable(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	container := "g6_itest_create_migration_table"

	out, err := exec.Command("docker", "run",
		"--rm",
		"--name", container,
		"--detach",
		"--publish", "5432:5432",
		"--env", "POSTGRES_USER=g6_test",
		"--env", "POSTGRES_DB=g6_test",
		"--env", "POSTGRES_PASSWORD=password",
		"postgres:alpine",
	).Output()
	assert.NoError(t, err, string(out))

	defer func(container string) {
		out, err := exec.Command("docker", "stop", container).Output()
		assert.NoError(t, err, string(out))
	}(container)
	for {
		time.Sleep(500 * time.Millisecond)
		out, err := exec.Command("docker", "logs", container, "--tail", "10").Output()
		output := string(out)
		fmt.Println(output)
		assert.NoError(t, err, output)

		if strings.Contains(output, "PostgreSQL init process complete; ready for start up") {
			break
		}

	}

	conn, err := sql.Open("postgres", "postgres://g6_test:password@127.0.0.1:5432/g6_test?sslmode=disable")
	assert.NoError(t, err)
	pg := &postgresSetup{conn}
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
			_, err := pg.CreateMigrationTable(tt.args.tableName)
			assert.Equal(t, fmt.Sprintf("%v", tt.expectedError), fmt.Sprintf("%v", err))
		})
	}
}
