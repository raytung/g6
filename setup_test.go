package g6

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/spf13/cobra"
	"github.com/raytung/g6/repositories"
	"database/sql"
)

func TestNewSetup(t *testing.T) {
	type args struct {
		migrationsRepo *mockMigrationsRepository
	}
	tests := []struct {
		name                          string
		args                          args
		cmd                           *cobra.Command
		expectedError                 error
		options                       *SetupOptions
		cmdArgs                       []string
		expectedCalledCreateTableWith []string
	}{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := NewSetup(tt.args.migrationsRepo)
			err := setup(tt.cmd, tt.cmdArgs, tt.options)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedCalledCreateTableWith, tt.args.migrationsRepo.calledCreateTableArgs)
		})
	}
}

type mockMigrationsRepository struct {
	result                sql.Result
	err                   error
	calledCreateTableArgs []string
}

var _ repositories.Migrations = &mockMigrationsRepository{}

func (m *mockMigrationsRepository) CreateTable(tableName string) (sql.Result, error) {
	m.calledCreateTableArgs = append(m.calledCreateTableArgs, tableName)
	return m.result, m.err
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
