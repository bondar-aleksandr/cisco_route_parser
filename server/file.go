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
	lock sync.Mutex
	Created bool
	Name string
	buffer     *bytes.Buffer
	OutputFile *os.File
}

func NewFile() *File {
	return &File{
		Name: uuid.NewString(),
		buffer: &bytes.Buffer{},
	}
}

func (f *File) SetFile() error {
	f.lock.Lock()
	file, err := os.Create(filepath.Join(tempFileDir, f.Name))
	if err != nil {
		return err
	}
	f.Created = true
	f.OutputFile = file
	return nil
}

func (f *File) Write(chunk []byte) error {
	if f.OutputFile == nil {
		return nil
	}
	_, err := f.OutputFile.Write(chunk)
	return err
}

func (f *File) Close() error {
	f.lock.Unlock()
	return f.OutputFile.Close()
}

func (f *File) Delete() error {
	return os.Remove(f.OutputFile.Name())
}