package g6

import (
	"github.com/spf13/cobra"
	"github.com/raytung/g6/repositories"
)

type CreateSetupService func(repositories.Migrations) SetupService
type SetupService func(*cobra.Command, []string, *SetupOptions) error

var _ CreateSetupService = NewSetup

type SetupOptions struct {
	table string
}

const (
	DefaultMigrationsTable = "g6_migrations"
)

func NewSetup(migrations repositories.Migrations) SetupService {
	return func(cmd *cobra.Command, args []string, options *SetupOptions) error {
		table := DefaultMigrationsTable
		if options != nil && options.table != ""{
			table = options.table
		}
		_, err := migrations.CreateTable(table)
		return err
	}
}
