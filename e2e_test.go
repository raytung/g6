package main

import (
	"testing"
	"os/exec"
)

func TestGenerate(t *testing.T) {
	testBinary := "g6_test"
	exec.Command("go", "build", "-o", testBinary).Output()
	defer func() {
		exec.Command("rm", testBinary).Run()
	}()
}