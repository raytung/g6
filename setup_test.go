package g6

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/raytung/g6/repositories"
	"database/sql"
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
	}
	tests := []testCase{
		{
			name: "creates migration table",
			args: args{&mockMigrationsRepository{
				err:    nil,
				result: &mockSQLResult{id: 1, rowsAffected: 1},
			}},
			options:                       nil,
			expectedError:                 nil,
			cmdArgs:                       []string{},
			expectedCalledCreateTableWith: []string{"g6_migrations"},
		},

		{
			name: "creates migration table (with table options)",
			args: args{&mockMigrationsRepository{
				err:                   nil,
				result:                &mockSQLResult{id: 1, rowsAffected: 1},
				calledCreateTableArgs: []string{},
			}},
			options:                       &SetupOptions{"other_migrations_table"},
			expectedError:                 nil,
			cmdArgs:                       []string{},
			expectedCalledCreateTableWith: []string{"other_migrations_table"},
		},

		{
			name: "creates migration table with default table name if options table is empty",
			args: args{&mockMigrationsRepository{
				err:                   nil,
				result:                &mockSQLResult{id: 1, rowsAffected: 1},
				calledCreateTableArgs: []string{},
			}},
			options:                       &SetupOptions{""},
			expectedError:                 nil,
			cmdArgs:                       []string{},
			expectedCalledCreateTableWith: []string{"g6_migrations"},
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
