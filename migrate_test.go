package g6

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/raytung/g6/repositories"
	"path/filepath"
	"time"
	"errors"
)

func TestNewMigrate(t *testing.T) {
	tests := []struct {
		name                           string
		args                           []string
		migrationsRepo                 *mockMigrationsRepo
		filePathReader                 *mockFilePathReader
		fileReader                     *mockFileReader
		expectedErr                    error
		options                        *MigrateOptions
		expectedCalledRunWithArgs      []*repositories.Migration
		expectedCalledReadFileWithArgs []string
		expectedCalledGlobWithArgs     []string
		expectedCallIsDirWithArgs      []string
	}{
		{
			name:        "No pending migrations",
			expectedErr: nil,
			migrationsRepo: &mockMigrationsRepo{
				latestMigration: &repositories.MigrationQueryResult{
					HasResults: true,
					ID:         1,
					Name:       "V1234_create_users_table",
					MigratedAt: time.Now(),
				},
				tableExist: true,
			},
			fileReader: &mockFileReader{
				isDir:             true,
				readFileResponses: [][]byte{[]byte("")},
			},
			filePathReader: &mockFilePathReader{
				files: []string{
					"V1234_create_users_table.up.sql",
				},
			},
			options:                    &MigrateOptions{"some_directory"},
			expectedCallIsDirWithArgs:  []string{"some_directory"},
			expectedCalledGlobWithArgs: []string{filepath.Join("some_directory", "*.up.sql")},
		},

		{
			name:        "No previous migrations",
			expectedErr: nil,
			migrationsRepo: &mockMigrationsRepo{
				latestMigration: &repositories.MigrationQueryResult{
					HasResults: false,
				},
				tableExist: true,
			},
			fileReader: &mockFileReader{
				readFileResponses: [][]byte{
					[]byte("CREATE TABLE users ();"),
				},
				readFileErrs: []error{nil, nil},
				isDir:        true,
			},
			filePathReader: &mockFilePathReader{
				files: []string{
					"V1234_create_users_table.up.sql",
				},
			},
			expectedCalledRunWithArgs: []*repositories.Migration{
				{Name: "V1234_create_users_table", Query: "CREATE TABLE users ();"},
			},
			expectedCalledReadFileWithArgs: []string{
				"V1234_create_users_table.up.sql",
			},
			options:                    &MigrateOptions{"some_directory"},
			expectedCalledGlobWithArgs: []string{filepath.Join("some_directory", "*.up.sql")},
			expectedCallIsDirWithArgs:  []string{"some_directory"},
		},

		{
			name:        "Some pending migrations",
			expectedErr: nil,
			migrationsRepo: &mockMigrationsRepo{
				latestMigration: &repositories.MigrationQueryResult{
					HasResults: true,
					ID:         1,
					Name:       "V1234_create_users_table",
					MigratedAt: time.Now(),
				},
				tableExist: true,
			},
			fileReader: &mockFileReader{
				readFileResponses: [][]byte{
					[]byte("CREATE TABLE posts ();"),
					[]byte("CREATE TABLE tags ();"),
				},
				readFileErrs: []error{nil, nil},
				isDir:        true,
			},
			filePathReader: &mockFilePathReader{
				files: []string{
					"V1234_create_users_table.up.sql",

					"V1235_create_posts_table.up.sql",

					"V1236_create_tags_table.up.sql",
				},
			},
			expectedCalledRunWithArgs: []*repositories.Migration{
				{Name: "V1235_create_posts_table", Query: "CREATE TABLE posts ();"},
				{Name: "V1236_create_tags_table", Query: "CREATE TABLE tags ();"},
			},
			expectedCalledReadFileWithArgs: []string{
				"V1235_create_posts_table.up.sql",
				"V1236_create_tags_table.up.sql",
			},
			options:                    &MigrateOptions{"some_directory"},
			expectedCalledGlobWithArgs: []string{filepath.Join("some_directory", "*.up.sql")},
			expectedCallIsDirWithArgs:  []string{"some_directory"},
		},

		{
			name:           "path is not a directory",
			expectedErr:    errors.New("not a directory"),
			fileReader:     &mockFileReader{isDir: false},
			options:        &MigrateOptions{"some directory"},
			migrationsRepo: &mockMigrationsRepo{},
			filePathReader: &mockFilePathReader{},
		},

		{
			name:           "error calling IsDir",
			expectedErr:    errors.New("some fs error"),
			fileReader:     &mockFileReader{isDirErr: errors.New("some fs error")},
			options:        &MigrateOptions{"some_directory"},
			migrationsRepo: &mockMigrationsRepo{},
			filePathReader: &mockFilePathReader{},
		},

		{
			name:                       "error calling Glob",
			expectedErr:                errors.New("some glob error"),
			fileReader:                 &mockFileReader{isDir: true},
			options:                    &MigrateOptions{"some_directory"},
			migrationsRepo:             &mockMigrationsRepo{},
			expectedCalledGlobWithArgs: []string{filepath.Join("some_directory", "*.up.sql")},
			filePathReader:             &mockFilePathReader{err: errors.New("some glob error")},
		},

		{
			name:                       "error getting latest migration",
			expectedErr:                errors.New("some db error"),
			fileReader:                 &mockFileReader{isDir: true},
			options:                    &MigrateOptions{"some_directory"},
			migrationsRepo:             &mockMigrationsRepo{latestErr: errors.New("some db error")},
			expectedCalledGlobWithArgs: []string{filepath.Join("some_directory", "*.up.sql")},
			filePathReader:             &mockFilePathReader{},
		},

		{
			name:        "error reading file",
			expectedErr: errors.New("some read file error"),
			fileReader: &mockFileReader{
				isDir:        true,
				readFileErrs: []error{errors.New("some read file error")},
				readFileResponses: [][]byte{
					[]byte("CREATE TABLE posts ();"),
				},
			},
			options: &MigrateOptions{"some_directory"},
			migrationsRepo: &mockMigrationsRepo{
				latestMigration: &repositories.MigrationQueryResult{
					HasResults: true,
					ID:         1,
					Name:       "V1234_create_users_table",
					MigratedAt: time.Now(),
				},
				tableExist: true,
			},
			expectedCalledGlobWithArgs:     []string{filepath.Join("some_directory", "*.up.sql")},
			expectedCalledReadFileWithArgs: []string{"V1235_create_posts_table.up.sql"},
			filePathReader: &mockFilePathReader{
				files: []string{
					"V1234_create_users_table.up.sql",
					"V1235_create_posts_table.up.sql",
				},
			},
		},

		{
			name:        "error running migrations",
			expectedErr: errors.New("some migration error"),
			fileReader: &mockFileReader{
				isDir:        true,
				readFileErrs: []error{nil},
				readFileResponses: [][]byte{
					[]byte("CREATE TABLE posts ();"),
				},
			},
			options: &MigrateOptions{"some_directory"},
			migrationsRepo: &mockMigrationsRepo{
				latestMigration: &repositories.MigrationQueryResult{
					HasResults: true,
					ID:         1,
					Name:       "V1234_create_users_table",
					MigratedAt: time.Now(),
				},
				runErr:     errors.New("some migration error"),
				tableExist: true,
			},
			expectedCalledGlobWithArgs:     []string{filepath.Join("some_directory", "*.up.sql")},
			expectedCalledReadFileWithArgs: []string{"V1235_create_posts_table.up.sql"},
			expectedCalledRunWithArgs: []*repositories.Migration{
				{Name: "V1235_create_posts_table", Query: "CREATE TABLE posts ();"},
			},
			filePathReader: &mockFilePathReader{
				files: []string{
					"V1234_create_users_table.up.sql",
					"V1235_create_posts_table.up.sql",
				},
			},
		},
	}
	for _, testCase := range tests {
		tt := testCase
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			migrate := NewMigrate(tt.migrationsRepo, tt.filePathReader, tt.fileReader)
			err := migrate(tt.args, tt.options)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedCalledRunWithArgs, tt.migrationsRepo.calledRunWithArgs, "Invalid args passed into Run")
			assert.Equal(t, tt.expectedCalledReadFileWithArgs, tt.fileReader.calledReadFileArgs, "Invalid args passed into ReadFile")
			assert.Equal(t, tt.expectedCalledGlobWithArgs, tt.filePathReader.calledGlobWithArgs, "Invalid args passed into Glob")
		})
	}
}

type mockFilePathReader struct {
	files              []string
	err                error
	calledGlobWithArgs []string
}

func (m *mockFilePathReader) Glob(pattern string) ([]string, error) {
	m.calledGlobWithArgs = append(m.calledGlobWithArgs, pattern)
	return m.files, m.err
}

var _ MigrationsRepository = &mockMigrationsRepo{}

type mockMigrationsRepo struct {
	latestMigration   *repositories.MigrationQueryResult
	latestErr         error
	tableExist        bool
	tableExistErr     error
	calledRunWithArgs []*repositories.Migration
	runErr            error
}

func (m *mockMigrationsRepo) Latest() (*repositories.MigrationQueryResult, error) {
	return m.latestMigration, m.latestErr
}

func (m *mockMigrationsRepo) TableExists() (bool, error) {
	return m.tableExist, m.tableExistErr
}

func (m *mockMigrationsRepo) Run(migration *repositories.Migration) error {
	m.calledRunWithArgs = append(m.calledRunWithArgs, migration)
	return m.runErr
}

type mockFileReader struct {
	calledReadFileArgs  []string
	readFileResponses   [][]byte
	readFileErrs        []error
	currReadFileCalls   int
	calledIsDirWithArgs []string
	isDir               bool
	isDirErr            error
}

func (m *mockFileReader) ReadFile(filename string) (content []byte, err error) {
	m.calledReadFileArgs = append(m.calledReadFileArgs, filename)
	content = m.readFileResponses[m.currReadFileCalls]
	err = m.readFileErrs[m.currReadFileCalls]
	m.currReadFileCalls += 1
	return
}

func (m *mockFileReader) IsDir(filename string) (bool, error) {
	m.calledIsDirWithArgs = append(m.calledIsDirWithArgs, filename)
	return m.isDir, m.isDirErr
}
