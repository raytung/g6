package g6

import (
	"github.com/spf13/cobra"
	"os"
	"time"
	"fmt"
)

type Opts interface {
	CreateFile(string) (*os.File, error)
	Mkdir(string) error
	PathJoin(...string) string
	TimeNow() time.Time
}

func NewGenerate(opts Opts) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "create SQL migration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := opts.Mkdir("migrations")
			if err != nil {
				if !os.IsExist(err) {
					return err
				}
			}
			now := opts.TimeNow().UnixNano() / int64(time.Millisecond)
			fileName := fmt.Sprintf("V%d__%s", now, args[0])
			path := opts.PathJoin("migrations", fileName)
			_, err = opts.CreateFile(path + ".up.sql")
			if err != nil {
				return err
			}
			_, err = opts.CreateFile(path + ".down.sql")
			return err
		},
	}
}
