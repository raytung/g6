package g6

import (
	"github.com/spf13/cobra"
	"github.com/raytung/g6/repositories"
)

type CreateSetupService func(repositories.Migrations) SetupService
type SetupService func(*cobra.Command, []string, *SetupOptions) error

var _ CreateSetupService = NewSetup

type SetupOptions struct {
}

func NewSetup(migrations repositories.Migrations) SetupService {
	return func(cmd *cobra.Command, args []string, options *SetupOptions) error {
		_, err := migrations.CreateTable("g6_migrations")
		return err
	}
}
