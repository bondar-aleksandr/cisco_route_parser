package main

import (
	"os"
	"path/filepath"
	"github.com/google/uuid"
)

const tempFileDir = "temp"

type File struct {
	filePath string
	Created bool
	Name string
	Platform string
	OutputFile *os.File
}

func NewFile() *File {
	return &File{
		Name: uuid.NewString(),
	}
}

func (f *File) SetFile(p string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	d := filepath.Join(cwd, tempFileDir)
	err = os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(d, f.Name))
	if err != nil {
		return err
	}
	f.filePath = filepath.Join(d, f.Name)
	f.Created = true
	f.OutputFile = file
	f.Platform = p
	return nil
}

func (f *File) Open() (*os.File, error) {
	r, err := os.Open(f.filePath)
	return r, err
}

func (f *File) Write(chunk []byte) error {
	_, err := f.OutputFile.Write(chunk)
	return err
}

func (f *File) Delete() error {
	err := f.OutputFile.Close()
	if err != nil {
		return err
	}
	return os.Remove(f.filePath)
}