package main

import (
	"testing"
	"os/exec"
	"strings"
	"os"
	"path/filepath"
	"github.com/stretchr/testify/assert"
	"database/sql"
	"github.com/raytung/g6/pkg/tests/docker"
)

const testBinary = "g6_test"

func TestE2EG6(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tmpDir := os.TempDir()
	binaryPath := filepath.Join(tmpDir, testBinary)

	defer tearDown(tmpDir)

	if err := exec.Command("go", "build", "-o", binaryPath).Run(); err != nil {
		t.Fatalf("%v", err)
	}

	t.Run("<no command>", func(t *testing.T) {
		output, err := exec.Command(binaryPath).Output()
		if err != nil {
			t.Fatalf("%v", err)
		}

		expectOutToContain := "g6 is a user friendly database migration tool"
		if !strings.Contains(string(output), expectOutToContain) {
			t.Fatalf("<no command> did not output string: %v", expectOutToContain)
		}
	})

	out, err, _ := docker.Cli(&docker.Options{
		Command:       "run",
		ContainerName: "test_e2e_g6",
		Image:         "postgres:alpine",
		Publish:       "5433:5432",
		Env: map[string]string{
			"POSTGRES_USER":     "g6_test",
			"POSTGRES_DB":       "g6_test",
			"POSTGRES_PASSWORD": "password",
		},
	})

	//defer dockerStop()
	assert.NoError(t, err, string(out))

	db, err := sql.Open("postgres", "postgres://g6_test:password@0.0.0.0:5433/g6_test?sslmode=disable")

	docker.WaitForDB(t, db)

	t.Run("setup", func(t *testing.T) {
		output, err := exec.Command(binaryPath,
			"setup",
			"--table", "g6_migrations",
			"--connection", "postgres://g6_test:password@0.0.0.0:5433/g6_test?sslmode=disable",
		).Output()
		assert.NoError(t, err)

		expectedOutput := ""
		assert.Equal(t, expectedOutput, string(output))
	})

	t.Run("migrate", func(t *testing.T) {
		output, err := exec.Command(binaryPath,
			"migrate",
			"--directory", "samples",
			"--table", "g6_migrations",
			"--connection", "postgres://g6_test:password@0.0.0.0:5433/g6_test?sslmode=disable",
		).Output()

		if assert.NoError(t, err) {
			expectedOutput := ""
			assert.Equal(t, expectedOutput, string(output))
		}

		_, err = db.Query("SELECT * FROM users LIMIT 1")
		assert.NoError(t, err)

		results, err := db.Query("SELECT * FROM g6_migrations WHERE name = $1", "V1524861198888__create_users_table")
		assert.NoError(t, err)

		assert.True(t, results.Next())

	})
}

func tearDown(path string) {
	os.RemoveAll(path)
	os.RemoveAll("migrations")
}
