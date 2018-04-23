package g6

import (
	"github.com/spf13/cobra"
)

func NewGenerateCmd(generateService GenerateService) *cobra.Command {
	genFlags := GenerateFlags{}
	cmd := &cobra.Command{
		Aliases: []string{"g"},
		Use:     "generate [file]",
		Example: "  g6 generate create_users_table\n  g6 generate create_posts_table --directory db/migrations",
		Short:   "create SQL migration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateService(cmd, args, &genFlags)
		},
	}

	cmd.Flags().StringVarP(&genFlags.directory, "directory", "d", "", "custom directory to look for SQL files");
	return cmd
}
