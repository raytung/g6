package g6

import (
	"github.com/spf13/cobra"
	"database/sql"
	"github.com/raytung/g6/repositories"
)

func NewSetupCmd() *cobra.Command {
	options := SetupOptions{}
	cmd := &cobra.Command{
		Use:     "setup",
		Example: `  g6 setup --table g6_migrations --connection "postgres://<username>:<password>@<host>:<port>/<db name>"`,
		Short:   "Setup database to keep track of migrations status",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := sql.Open("postgres", options.dbConnection)
			if err != nil {
				return err
			}
			migrationsRepo := repositories.NewPostgresMigrations(db)
			setupService := NewSetup(migrationsRepo)
			return setupService(args, &options)
		},
	}

	cmd.Flags().StringVarP(&options.table, "table", "t", "", "g6 migrations table");
	cmd.Flags().StringVarP(&options.dbConnection, "connection", "c", "", "connection string");
	return cmd
}
