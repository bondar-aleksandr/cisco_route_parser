package main

import (
	"bytes"
	"os"
	"sync"
	// "log"
	"path/filepath"
	"github.com/google/uuid"
)

const tempFileDir = "./temp/"

type File struct {
	filePath string
	lock sync.RWMutex
	Created bool
	Name string
	Platform string
	buffer     *bytes.Buffer
	OutputFile *os.File
}

func NewFile() *File {
	return &File{
		Name: uuid.NewString(),
		buffer: &bytes.Buffer{},
	}
}

func (f *File) SetFile(p string) error {
	f.filePath = filepath.Join(tempFileDir, f.Name)
	file, err := os.Create(f.filePath)
	if err != nil {
		return err
	}
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
	if f.OutputFile == nil {
		return nil
	}
	f.lock.Lock()
	_, err := f.OutputFile.Write(chunk)
	f.lock.Unlock()
	return err
}

func (f *File) Close() error {
	return f.OutputFile.Close()
}

func (f *File) Delete() error {
	return os.Remove(f.filePath)
}