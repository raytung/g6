package services

import "path/filepath"

type FilePath interface {
	Join(...string) string
}

type filePath struct {
}

func (f *filePath) Join(paths ...string) string {
	return filepath.Join(paths...)
}

func NewFilePathService() *filePath {
	return &filePath{}
}
