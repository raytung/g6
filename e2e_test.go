package main

import (
	"testing"
	"os/exec"
	"strings"
)

const testBinary = "g6_test"

func TestG6(t *testing.T) {
	defer tearDown()
	if err := exec.Command("go", "build", "-o", testBinary).Run(); err != nil {
		t.Fatalf("%v", err)
	}

	t.Run("<no command>", func(t *testing.T) {
		output, err := exec.Command("./" + testBinary).Output()
		if err != nil {
			t.Fatalf("%v", err)
		}

		expectOutToContain := "g6 is a user friendly database migration tool"
		if !strings.Contains(string(output), expectOutToContain) {
			t.Fatalf("<no command> did not output string: %v", expectOutToContain)
		}
	})

}

func tearDown() {
	exec.Command("rm", testBinary).Run()
}