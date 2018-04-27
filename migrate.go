package g6

import (
	"github.com/raytung/g6/repositories"
	"path/filepath"
	"sort"
	"strings"
	"errors"
)

type CreateMigrateService func(MigrationsRepository, filePathReader, fileReader) MigrateService
type MigrateService func([]string, *MigrateOptions) error

var _ CreateMigrateService = NewMigrate

type MigrateOptions struct {
	directory string
}

type MigrationsRepository interface {
	repositories.MigrationsLatestInfo
	repositories.MigrationsRunner
}

type filePathReader interface {
	// Glob - Globs file path and return all matched file names
	Glob(pattern string) (matches []string, err error)
}

type fileReader interface {
	ReadFile(filename string) ([]byte, error)
	IsDir(filename string) (bool, error)
}

func NewMigrate(migrations MigrationsRepository, filePathReader filePathReader, fileReader fileReader) MigrateService {
	return func(args []string, options *MigrateOptions) error {
		if isDir, _ := fileReader.IsDir(options.directory); !isDir {
			return errors.New("not a directory")
		}
		upFiles, _ := filePathReader.Glob(filepath.Join(options.directory, "*.up.sql"))

		sort.Sort(sort.StringSlice(upFiles))

		latest, _ := migrations.Latest()

		var latestMigrationIndex int

		for index, fileName := range upFiles {
			if strings.Contains(fileName, latest.Name) {
				latestMigrationIndex = index
			}
		}

		if len(upFiles) == latestMigrationIndex {
			return nil
		}

		pendingMigrations := upFiles[latestMigrationIndex+1:]

		for _, fileName := range pendingMigrations {
			name := fileName[:len(fileName) - len(".up.sql")]
			m := repositories.Migration{Name: name}
			content, _ := fileReader.ReadFile(fileName)
			m.Query = string(content)
			migrations.Run(&m)
		}
		return nil
	}
}
