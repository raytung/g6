package g6

import (
	"github.com/spf13/cobra"
	"database/sql"
	"github.com/raytung/g6/repositories"
)

const (
	TableFlag      = "table"
	TableShortFlag = "t"

	ConnectionStringFlag      = "connection"
	ConnectionStringShortFlag = "c"
)

func NewSetupCmd() *cobra.Command {
	options := SetupOptions{}
	cmd := &cobra.Command{
		Use:     "setup",
		Example: `  g6 setup --table g6_migrations --connection "postgres://<username>:<password>@<host>:<port>/<db name>"`,
		Short:   "Setup database to keep track of migrationsRepo status",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := sql.Open("postgres", cmd.Flag(ConnectionStringFlag).Value.String())
			if err != nil {
				return err
			}
			migrationsRepo := repositories.NewPostgresMigrations(db, cmd.Flag(TableFlag).Value.String())
			setupService := NewSetup(migrationsRepo)
			return setupService(args, &options)
		},
	}

	cmd.Flags().StringP(TableFlag, TableShortFlag, DefaultMigrationsTable, "g6 migrationsRepo table")
	cmd.Flags().StringP(ConnectionStringFlag, ConnectionStringShortFlag, "", "connection string")
	return cmd
}
