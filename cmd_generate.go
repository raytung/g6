package g6

import (
	"github.com/spf13/cobra"
)

func NewGenerateCmd(generateService GenerateService) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "create SQL migration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateService(cmd, args)
		},
	}
}
