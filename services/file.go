package services

import "os"

type File interface {
	Create(string) (*os.File, error)
	Mkdir(string) error
}

type file struct {
}

func (f *file) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (f *file) Mkdir(name string) error {
	return os.Mkdir(name, os.ModePerm)
}

func NewFileService() File {
	return &file{}
}