package g6

import (
	"github.com/spf13/cobra"
)

func NewGenerateCmd(generateService GenerateService) *cobra.Command {
	genFlags := GenerateFlags{}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "create SQL migration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateService(cmd, args, &genFlags)
		},
	}

	cmd.Flags().StringVarP(&genFlags.directory, "directory", "d", "", "custom directory to look for SQL files");
	return cmd
}
