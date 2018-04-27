package g6

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/raytung/g6/repositories"
	"database/sql"
	"errors"
)

func TestNewSetup(t *testing.T) {
	type args struct {
		migrationsRepo *mockMigrationsRepository
	}
	type testCase struct {
		name          string
		args          args
		expectedError error
		options       *SetupOptions
		cmdArgs       []string
	}
	tests := []testCase{
		{
			name: "creates migration table",
			args: args{&mockMigrationsRepository{
				result: &mockSQLResult{id: 1, rowsAffected: 1},
			}},
		},

		{
			name: "creates migration table (with table options)",
			args: args{&mockMigrationsRepository{
				result: &mockSQLResult{id: 1, rowsAffected: 1},
			}},
			options: &SetupOptions{table: "other_migrations_table"},
		},

		{
			name: "creates migration table with default table name if options table is empty",
			args: args{&mockMigrationsRepository{
				result: &mockSQLResult{id: 1, rowsAffected: 1},
			}},
			options: &SetupOptions{table: ""},
		},

		{
			name: "does not attempt to create migration table if it already exists",
			args: args{&mockMigrationsRepository{
				tableExist: true,
				result:     &mockSQLResult{id: 1, rowsAffected: 1},
			}},
			options: &SetupOptions{table: ""},
		},

		{
			name: "does not attempt to create migration table TableExists errors out",
			args: args{&mockMigrationsRepository{
				tableExist:    true,
				tableExistErr: errors.New("some table exist error"),
				result:        &mockSQLResult{id: 1, rowsAffected: 1},
			}},
			options:       &SetupOptions{table: ""},
			expectedError: errors.New("some table exist error"),
		},
	}
	for _, tt := range tests {
		testCase := tt // copy to prevent racing
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			setup := NewSetup(testCase.args.migrationsRepo)
			err := setup(testCase.cmdArgs, testCase.options)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

type mockMigrationsRepository struct {
	result                sql.Result
	err                   error
	calledCreateTableArgs []string
	calledTableExistsArgs []string
	tableExist            bool
	tableExistErr         error
}

var _ repositories.Migrations = &mockMigrationsRepository{}

func (m *mockMigrationsRepository) CreateTable() (sql.Result, error) {
	return m.result, m.err
}

func (m *mockMigrationsRepository) TableExists() (bool, error) {
	return m.tableExist, m.tableExistErr
}

func (m *mockMigrationsRepository) Latest() (*repositories.MigrationQueryResult, error) {
	return nil, nil
}

type mockSQLResult struct {
	id              int64
	rowsAffected    int64
	err             error
	rowsAffectedErr error
}

var _ sql.Result = &mockSQLResult{}

func (m *mockSQLResult) LastInsertId() (int64, error) {
	return m.id, m.err
}

func (m *mockSQLResult) RowsAffected() (int64, error) {
	return m.rowsAffected, m.rowsAffectedErr
}
