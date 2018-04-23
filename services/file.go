package services

import "os"

type File interface {
	Create(string) (*os.File, error)
	Mkdir(string) error
	IsExist(error) bool
}

type file struct {
}

func (f *file) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (f *file) Mkdir(name string) error {
	return os.MkdirAll(name, os.ModePerm)
}

func (f *file) IsExist(err error) bool {
	return os.IsExist(err)
}

func NewFileService() File {
	return &file{}
}
