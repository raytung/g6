package g6

import (
	"os"
	"time"
	"github.com/spf13/cobra"
	"fmt"
	"github.com/raytung/g6/services"
)

type CreateGenerateService func(services.File, services.FilePath, services.Time) GenerateService
type GenerateService func(*cobra.Command, []string) error

var _ CreateGenerateService = NewGenerate

func NewGenerate(file services.File, path services.FilePath, time2 services.Time) GenerateService {
	return func(cmd *cobra.Command, args []string) error {
		err := file.Mkdir("migrations")
		if err != nil {
			if !os.IsExist(err) {
				return err
			}
		}
		now := time2.TimeNow().UnixNano() / int64(time.Millisecond)
		fileName := fmt.Sprintf("V%d__%s", now, args[0])
		path := path.Join("migrations", fileName)
		_, err = file.Create(path + ".up.sql")
		if err != nil {
			return err
		}
		_, err = file.Create(path + ".down.sql")
		return err
	}
}
