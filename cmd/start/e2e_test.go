package main

import (
	"testing"
	"os/exec"
	"strings"
	"os"
	"path/filepath"
	"fmt"
)

const testBinary = "g6_test"

func TestG6(t *testing.T) {
	tmpDir := os.TempDir()
	binaryPath := filepath.Join(tmpDir, testBinary)

	defer tearDown(binaryPath)

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

	t.Run("generate", func(t *testing.T) {
		out, err := exec.Command(binaryPath, "generate", "create_users_table").Output()
		fmt.Println(string(out))
		if err != nil {
			t.Fatalf("%v", err)
		}

		path := filepath.Join("migrations", "V*__create_users_table.*.sql")
		files, err := filepath.Glob(path)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if len(files) != 2 {
			t.Fatalf("Generated more than 2 sql file: %v", files)
		}

		for _, file := range files {
			fileInfo, err := os.Stat(file)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if fileInfo.IsDir() {
				t.Fatalf("Is a directory")
			}
		}
	})
}

func tearDown(path string) {
	exec.Command("rm", path).Run()
	os.RemoveAll("migrations")
}
