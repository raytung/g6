package main

import (
	"testing"
	"os/exec"
	"os"
	"path/filepath"
	"github.com/stretchr/testify/assert"
	"fmt"
)

const testBinaryForGenerate = "g6_gen__test"

func TestGenerate(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tmpDir := os.TempDir()
	binaryPath := filepath.Join(tmpDir, testBinaryForGenerate)

	defer tearDown(tmpDir)

	err := exec.Command("go", "build", "-o", binaryPath).Run()
	assert.NoError(t, err)

	testCases := []struct {
		name              string
		expectedFilesGlob string
		cmd               []string
	}{
		{
			name:              "no flags",
			expectedFilesGlob: filepath.Join("migrations", "V*__create_users_table.*.sql"),
			cmd:               []string{"generate", "create_users_table"},
		},
		{
			name:              "directory long flag (--directory)",
			expectedFilesGlob: filepath.Join(tmpDir, "db", "migrations", "V*__create_users_table.*.sql"),
			cmd:               []string{"generate", "create_users_table", "--directory", filepath.Join(tmpDir, "db", "migrations")},
		},
		{
			name:              "directory short flag (-d)",
			expectedFilesGlob: filepath.Join(tmpDir, "new_migrations", "V*__create_blogs_table.*.sql"),
			cmd:               []string{"generate", "create_blogs_table", "-d", filepath.Join(tmpDir, "new_migrations")},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			out, err := exec.Command(binaryPath, testCase.cmd...).Output()

			fmt.Println(string(out))
			assert.NoError(t, err)

			files, err := filepath.Glob(testCase.expectedFilesGlob)
			assert.NoError(t, err)

			if len(files) != 2 {
				t.Fatalf("Did not generate 2 sql files: %v", files)
			}

			for _, file := range files {
				fileInfo, err := os.Stat(file)
				assert.NoError(t, err)
				assert.False(t, fileInfo.IsDir(), "Is a directory")
			}
		})
	}
}
