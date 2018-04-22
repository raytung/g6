package g6

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"path/filepath"
	"time"
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

type opts struct {
}

func (o *opts) CreateFile(name string) (*os.File, error) {
	return os.Create(name)
}

func (o *opts) PathJoin(path ...string) string {
	return filepath.Join(path...)
}

func (o *opts) Mkdir(name string) error {
	return os.Mkdir(name, os.ModePerm)
}

func (o *opts) TimeNow() time.Time {
	return time.Now()
}

func Execute() {
	optsService := opts{}
	rootCmd := NewRoot()
	rootCmd.AddCommand(NewGenerate(&optsService))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
