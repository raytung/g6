package g6

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"github.com/raytung/g6/services"
)

func NewRoot() *cobra.Command {
	return &cobra.Command{
		Use:   "g6",
		Short: "g6 is a user friendly database migration tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}

func Execute() {
	fileService := services.NewFileService()
	timeService := services.NewTimeService()
	generateService := NewGenerate(fileService, timeService)
	generateCmd := NewGenerateCmd(generateService)
	rootCmd := NewRoot()
	rootCmd.AddCommand(generateCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
