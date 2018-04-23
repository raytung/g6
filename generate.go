package g6

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/raytung/g6/services"
	"path/filepath"
)

type CreateGenerateService func(services.File, services.VersionGenerator) GenerateService
type GenerateService func(*cobra.Command, []string) error

var _ CreateGenerateService = NewGenerate

func NewGenerate(file services.File, versionGen services.VersionGenerator) GenerateService {
	return func(cmd *cobra.Command, args []string) error {
		err := file.Mkdir("migrations")
		if err != nil {
			if !file.IsExist(err) {
				return err
			}
		}
		fileName := fmt.Sprintf("V%s__%s", versionGen.Generate(), args[0])
		path := filepath.Join("migrations", fileName)
		_, err = file.Create(path + ".up.sql")
		if err != nil {
			return err
		}
		_, err = file.Create(path + ".down.sql")
		return err
	}
}
