package g6

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/raytung/g6/services"
	"path/filepath"
	"strings"
)

type CreateGenerateService func(services.File, services.VersionGenerator) GenerateService
type GenerateService func(*cobra.Command, []string, *GenerateFlags) error

var _ CreateGenerateService = NewGenerate

type GenerateFlags struct {
	directory string
}

func NewGenerate(file services.File, versionGen services.VersionGenerator) GenerateService {
	return func(cmd *cobra.Command, args []string, genFlags *GenerateFlags) error {
		migrationsDir := "migrations"
		if genFlags != nil {
			if strings.TrimSpace(genFlags.directory) != "" {
				migrationsDir = genFlags.directory
			}
		}
		err := file.Mkdir(migrationsDir)
		if err != nil {
			if !file.IsExist(err) {
				return err
			}
		}

		fileName := fmt.Sprintf("V%s__%s", versionGen.Generate(), args[0])
		path := filepath.Join(migrationsDir, fileName)
		_, err = file.Create(path + ".up.sql")
		if err != nil {
			return err
		}
		_, err = file.Create(path + ".down.sql")
		return err
	}
}
