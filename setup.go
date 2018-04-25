package g6

import (
	"github.com/raytung/g6/repositories"
)

type CreateSetupService func(repositories.Migrations) SetupService
type SetupService func([]string, *SetupOptions) error

var _ CreateSetupService = NewSetup

type SetupOptions struct {
	table string
}

const (
	DefaultMigrationsTable = "g6_migrations"
)

func NewSetup(migrations repositories.Migrations) SetupService {
	return func(args []string, options *SetupOptions) error {
		table := DefaultMigrationsTable
		if options != nil && options.table != "" {
			table = options.table
		}
		if exists, err := migrations.TableExists(table); exists || err != nil {
			return err
		}
		_, err := migrations.CreateTable(table)
		return err
	}
}
