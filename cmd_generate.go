package g6

import (
	"github.com/spf13/cobra"
)

const (
	DirectoryFlag      = "directory"
	DirectoryShortFlag = "d"
)

func NewGenerateCmd(generateService GenerateService) *cobra.Command {
	cmd := &cobra.Command{
		Aliases: []string{"g"},
		Use:     "generate [file]",
		Example: "  g6 generate create_users_table\n  g6 generate create_posts_table --directory db/migrations",
		Short:   "create SQL migration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateService(args, &GenerateOptions{
				directory: cmd.Flag(DirectoryFlag).Value.String(),
			})
		},
	}

	cmd.Flags().StringP(DirectoryFlag, DirectoryShortFlag, "migrations", "custom directory to look for SQL files");
	return cmd
}
