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
		name                          string
		args                          args
		expectedError                 error
		options                       *SetupOptions
		cmdArgs                       []string
		expectedCalledCreateTableWith []string
		expectedCalledTableExistsWith []string
	}
	tests := []testCase{
		{
			name: "creates migration table",
			args: args{&mockMigrationsRepository{
				result: &mockSQLResult{id: 1, rowsAffected: 1},
			}},
			expectedCalledCreateTableWith: []string{"g6_migrations"},
			expectedCalledTableExistsWith: []string{"g6_migrations"},
		},

		{
			name: "creates migration table (with table options)",
			args: args{&mockMigrationsRepository{
				result: &mockSQLResult{id: 1, rowsAffected: 1},
			}},
			options:                       &SetupOptions{table: "other_migrations_table"},
			expectedCalledCreateTableWith: []string{"other_migrations_table"},
			expectedCalledTableExistsWith: []string{"other_migrations_table"},
		},

		{
			name: "creates migration table with default table name if options table is empty",
			args: args{&mockMigrationsRepository{
				result: &mockSQLResult{id: 1, rowsAffected: 1},
			}},
			options:                       &SetupOptions{table: ""},
			expectedCalledCreateTableWith: []string{"g6_migrations"},
			expectedCalledTableExistsWith: []string{"g6_migrations"},
		},

		{
			name: "does not attempt to create migration table if it already exists",
			args: args{&mockMigrationsRepository{
				tableExist: true,
				result:     &mockSQLResult{id: 1, rowsAffected: 1},
			}},
			options:                       &SetupOptions{table: ""},
			expectedCalledTableExistsWith: []string{"g6_migrations"},
		},

		{
			name: "does not attempt to create migration table TableExists errors out",
			args: args{&mockMigrationsRepository{
				tableExist:    true,
				tableExistErr: errors.New("some table exist error"),
				result:        &mockSQLResult{id: 1, rowsAffected: 1},
			}},
			options:                       &SetupOptions{table: ""},
			expectedError:                 errors.New("some table exist error"),
			expectedCalledTableExistsWith: []string{"g6_migrations"},
		},
	}
	for _, tt := range tests {
		testCase := tt // copy to prevent racing
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			setup := NewSetup(testCase.args.migrationsRepo)
			err := setup(testCase.cmdArgs, testCase.options)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedCalledCreateTableWith, testCase.args.migrationsRepo.calledCreateTableArgs)
			assert.Equal(t, testCase.expectedCalledTableExistsWith, testCase.args.migrationsRepo.calledTableExistsArgs)
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

func (m *mockMigrationsRepository) CreateTable(tableName string) (sql.Result, error) {
	m.calledCreateTableArgs = append(m.calledCreateTableArgs, tableName)
	return m.result, m.err
}

func (m *mockMigrationsRepository) TableExists(tableName string) (bool, error) {
	m.calledTableExistsArgs = append(m.calledTableExistsArgs, tableName)
	return m.tableExist, m.tableExistErr
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
