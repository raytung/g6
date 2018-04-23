package g6

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/raytung/g6/services"
	"path/filepath"
	"strings"
	"errors"
)

type CreateGenerateService func(services.File, services.VersionGenerator) GenerateService
type GenerateService func(*cobra.Command, []string, *GenerateFlags) error

var _ CreateGenerateService = NewGenerate

type GenerateFlags struct {
	directory string
}

const (
	DefaultMigrationsDirectory = "migrations"
	SqlPostfix                 = ".sql"
	UpFilePostfix              = ".up" + SqlPostfix
	DownFilePostfix            = ".down" + SqlPostfix
)

func NewGenerate(file services.File, versionGen services.VersionGenerator) GenerateService {
	return func(cmd *cobra.Command, args []string, genFlags *GenerateFlags) error {
		if len(args) == 0 {
			return errors.New("must provide migration file name")
		}

		dir := migrationDir(genFlags)

		if err := file.Mkdir(dir); err != nil && !file.IsExist(err) {
			return err
		}

		path := fullFilePath(versionGen.Generate(), dir, args[0])

		_, err := file.Create(path + UpFilePostfix)
		if err != nil {
			return err
		}
		_, err = file.Create(path + DownFilePostfix)
		return err
	}
}

func fullFilePath(version, directory, fileName string) string {
	strippedFileName := stripSqlPostfix(fileName)
	fullFileName := fmt.Sprintf("V%s__%s", version, strippedFileName)
	return filepath.Join(directory, fullFileName)
}

func stripSqlPostfix(path string) string {
	newPath := path
	if strings.HasSuffix(newPath, SqlPostfix) {
		strIndex := strings.LastIndex(newPath, SqlPostfix)
		newPath = newPath[0:strIndex]
	}
	return newPath
}

func migrationDir(genFlags *GenerateFlags) string {
	migrationsDir := DefaultMigrationsDirectory
	if genFlags != nil && strings.TrimSpace(genFlags.directory) != "" {
		migrationsDir = genFlags.directory
	}
	return migrationsDir
}
