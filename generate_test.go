package g6

import (
	"testing"

	"github.com/spf13/cobra"
	"os"
	"time"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestNewGenerate(t *testing.T) {
	type svcs struct {
		file  *fileSvcMock
		path  *filePathSvcMock
		time2 *timeSvcMock
	}
	tests := []struct {
		name          string
		services      svcs
		wantErr       bool
		want          GenerateService
		cmd           *cobra.Command
		cmdArgs       []string
		expectedError error
	}{
		{
			name: "happy path",
			services: svcs{
				file: &fileSvcMock{},
				path: &filePathSvcMock{
					path: "migrations/V1234__create_users_table",
				},
				time2: &timeSvcMock{},
			},
			wantErr: false,
			cmd:     nil,
			cmdArgs: []string{"create_users_table"},
		},

		{
			name: "Directory exists",
			services: svcs{
				file: &fileSvcMock{
					isExist:  true,
					mkdirErr: errors.New("some error"),
				},
				path: &filePathSvcMock{
					path: "migrations/V1234__create_users_table",
				},
				time2: &timeSvcMock{},
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
				path: &filePathSvcMock{
					path: "migrations/V1234__create_users_table",
				},
				time2: &timeSvcMock{},
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
				path: &filePathSvcMock{
					path: "migrations/V1234__create_users_table",
				},
				time2: &timeSvcMock{},
			},
			wantErr: true,
			cmd:     nil,
			cmdArgs: []string{"create_users_table"},
			expectedError: errors.New("some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewGenerate(tt.services.file, tt.services.path, tt.services.time2)
			err := gen(tt.cmd, tt.cmdArgs)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, tt.services.file.calledCreateArgs, "migrations/V1234__create_users_table.up.sql")
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

type filePathSvcMock struct {
	path string
}

func (f *filePathSvcMock) Join(paths ...string) string {
	return f.path
}

type timeSvcMock struct {
	now time.Time
}

func (t *timeSvcMock) TimeNow() time.Time {
	return t.now
}
