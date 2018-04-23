package main

import (
	"testing"
	"os/exec"
	"strings"
	"os"
	"path/filepath"
)

const testBinary = "g6_test"

func TestG6(t *testing.T) {
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
}

func tearDown(path string) {
	os.RemoveAll(path)
	os.RemoveAll("migrations")
}
