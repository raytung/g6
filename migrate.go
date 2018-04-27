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
		isDir, err := fileReader.IsDir(options.directory)
		if err != nil {
			return err
		}
		if !isDir {
			return errors.New("not a directory")
		}
		upFiles, err := filePathReader.Glob(filepath.Join(options.directory, "*.up.sql"))

		if err != nil {
			return err
		}

		sort.Sort(sort.StringSlice(upFiles))

		latest, err := migrations.Latest()

		if err != nil {
			return err
		}

		var latestMigrationIndex int

		for index, fileName := range upFiles {
			if strings.Contains(fileName, latest.Name) {
				latestMigrationIndex = index
			}
		}

		if latest.HasResults && len(upFiles)-1 == latestMigrationIndex {
			return nil
		}

		pendingMigrations := upFiles[latestMigrationIndex+1:]
		if !latest.HasResults {
			pendingMigrations = upFiles
		}

		for _, fileName := range pendingMigrations {
			name := fileName[:len(fileName)-len(".up.sql")]
			content, err := fileReader.ReadFile(fileName)

			if err != nil {
				return err
			}

			err = migrations.Run(&repositories.Migration{
				Name:  name,
				Query: string(content),
			})

			if err != nil {
				return err
			}
		}
		return nil
	}
}
