package g6

import (
	"testing"

	"github.com/spf13/cobra"
	"os"
	"github.com/stretchr/testify/assert"
	"errors"
	"path/filepath"
)

func TestNewGenerate(t *testing.T) {
	type svcs struct {
		file    *fileSvcMock
		version *versionSvcMock
	}
	tests := []struct {
		name               string
		services           svcs
		wantErr            bool
		want               GenerateService
		cmd                *cobra.Command
		cmdArgs            []string
		expectedError      error
		expectCreatedFiles []string
		genFlags           *GenerateFlags
	}{
		{
			name: "happy path",
			services: svcs{
				file: &fileSvcMock{},
				version: &versionSvcMock{
					gen: "0001",
				},
			},
			wantErr: false,
			cmd:     nil,
			cmdArgs: []string{"create_users_table"},
			expectCreatedFiles: []string{
				filepath.Join("migrations", "V0001__create_users_table.up.sql"),
				filepath.Join("migrations", "V0001__create_users_table.down.sql"),
			},
		},

		{
			name: "happy path (with directory flag but arg is empty)",
			services: svcs{
				file: &fileSvcMock{},
				version: &versionSvcMock{
					gen: "0001",
				},
			},
			wantErr: false,
			cmd:     nil,
			cmdArgs: []string{"create_users_table"},
			expectCreatedFiles: []string{
				filepath.Join("migrations", "V0001__create_users_table.up.sql"),
				filepath.Join("migrations", "V0001__create_users_table.down.sql"),
			},
			genFlags: &GenerateFlags{""},
		},

		{
			name: "happy path (with directory flag)",
			services: svcs{
				file: &fileSvcMock{},
				version: &versionSvcMock{
					gen: "0001",
				},
			},
			wantErr: false,
			cmd:     nil,
			cmdArgs: []string{"create_users_table"},
			expectCreatedFiles: []string{
				filepath.Join("new_migrations_dir", "V0001__create_users_table.up.sql"),
				filepath.Join("new_migrations_dir", "V0001__create_users_table.down.sql"),
			},
			genFlags: &GenerateFlags{"new_migrations_dir"},
		},

		{
			name: "No args",
			services: svcs{
				file:    &fileSvcMock{},
				version: &versionSvcMock{"0001"},
			},
			wantErr:       true,
			cmd:           nil,
			cmdArgs:       []string{},
			genFlags:      &GenerateFlags{"new_migrations_dir"},
			expectedError: errors.New("must provide migration file name"),
		},

		{
			name: "Directory exists",
			services: svcs{
				file: &fileSvcMock{
					isExist:  true,
					mkdirErr: errors.New("some error"),
				},
				version: &versionSvcMock{"0002"},
			},
			expectCreatedFiles: []string{
				filepath.Join("migrations", "V0002__create_users_table.up.sql"),
				filepath.Join("migrations", "V0002__create_users_table.down.sql"),
			},
			wantErr: false,
			cmd:     nil,
			cmdArgs: []string{"create_users_table"},
		},

		{
			name: "Unknown error while creating directory",
			services: svcs{
				file: &fileSvcMock{
					isExist:  false,
					mkdirErr: errors.New("some error"),
				},
				version: &versionSvcMock{},
			},
			wantErr:       true,
			cmd:           nil,
			cmdArgs:       []string{"create_users_table"},
			expectedError: errors.New("some error"),
		},

		{
			name: "Error while creating up.sql",
			services: svcs{
				file: &fileSvcMock{
					mkdirErr:  nil,
					createErr: errors.New("some error"),
				},
				version: &versionSvcMock{},
			},
			wantErr:       true,
			cmd:           nil,
			cmdArgs:       []string{"create_users_table"},
			expectedError: errors.New("some error"),
		},

		{
			name: "Strip .sql if name already has the same postfix",
			services: svcs{
				file: &fileSvcMock{
					mkdirErr:  nil,
					createErr: nil,
				},
				version: &versionSvcMock{"0001"},
			},
			wantErr: false,
			cmd:     nil,
			cmdArgs: []string{"create_users_table.sql"},
			expectCreatedFiles: []string{
				filepath.Join("migrations", "V0001__create_users_table.up.sql"),
				filepath.Join("migrations", "V0001__create_users_table.down.sql"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewGenerate(tt.services.file, tt.services.version)
			err := gen(tt.cmd, tt.cmdArgs, tt.genFlags)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.services.file.calledCreateArgs, tt.expectCreatedFiles)
			}
		})
	}
}

type fileSvcMock struct {
	isExist          bool
	createErr        error
	mkdirErr         error
	calledCreateArgs []string
	calledMkdirArgs  []string
}

func (f *fileSvcMock) Create(path string) (*os.File, error) {
	f.calledCreateArgs = append(f.calledCreateArgs, path)
	return nil, f.createErr
}

func (f *fileSvcMock) Mkdir(dir string) error {
	f.calledMkdirArgs = append(f.calledMkdirArgs, dir)
	return f.mkdirErr
}

func (f *fileSvcMock) IsExist(err error) bool {
	return f.isExist
}

type versionSvcMock struct {
	gen string
}

func (t *versionSvcMock) Generate() string {
	return t.gen
}
