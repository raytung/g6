package tests

import (
	"os/exec"
	"testing"
	"database/sql"
	"time"
	"fmt"
)

type DockerOptions struct {
	Command       string
	ContainerName string
	Publish       string
	Env           map[string]string
	Image         string
	Args          []string
}

type TearDown func() ([]byte, error)

func DockerCli(options *DockerOptions) ([]byte, error, TearDown) {
	args := []string{
		options.Command,
		"--rm",
		"--name", options.ContainerName,
		"--detach",
		"--publish", options.Publish,
	}
	for key, value := range options.Env {
		args = append(args, "--env", key+"="+value)
	}
	args = append(args, options.Image)
	for _, arg := range options.Args {
		args = append(args, arg)
	}
	out, err := exec.Command("docker", args...).Output()

	tearDown := func() ([]byte, error) {
		return exec.Command("docker", "stop", options.ContainerName).Output()
	}

	return out, err, tearDown
}

func WaitForDB(t *testing.T, db *sql.DB) {
	t.Helper()
	timeout := 30 * time.Second
	elapsed := 0 * time.Second
	waitTime := 500 * time.Millisecond

	for {
		time.Sleep(waitTime)
		elapsed += waitTime

		if elapsed >= timeout {
			t.Fatalf("Cannot connect to database after %v", timeout)
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
