package g6

import (
	"github.com/spf13/cobra"
	"github.com/raytung/g6/repositories"
	"database/sql"
	"path/filepath"
	"io/ioutil"
	"os"
)

type filePath struct {
}

func (f *filePath) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

type filereader struct {
}

func (f *filereader) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (f *filereader) IsDir(filename string) (bool, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func NewMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate [file]",
		Example: `  g6 migrate --directory="db/migrations" --connection="postgres://<username>:<password>@<host>:<port>/<db name>"`,
		Short:   "run pending migrations against database",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := sql.Open("postgres", cmd.Flag(ConnectionStringFlag).Value.String())
			if err != nil {
				return err
			}
			migrationsRepo := repositories.NewPostgresMigrations(db, cmd.Flag(TableFlag).Value.String())

			migrateService := NewMigrate(migrationsRepo, &filePath{}, &filereader{})

			return migrateService(args, &MigrateOptions{
				directory: cmd.Flag(DirectoryFlag).Value.String(),
			})
		},
	}

	cmd.Flags().StringP(DirectoryFlag, DirectoryShortFlag, "migrations", "custom directory to look for SQL files");
	cmd.Flags().StringP(ConnectionStringFlag, ConnectionStringShortFlag, "", "connection string")
	cmd.Flags().StringP(TableFlag, TableShortFlag, DefaultMigrationsTable, "g6 migrationsRepo table")

	return cmd
}
