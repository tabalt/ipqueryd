package pidfile

import (
	"io/ioutil"
	"os"
)

type File struct {
	Path string
}

func NewFile(path string) *File {
	return &File{
		Path: path,
	}
}

func (f *File) Read() ([]byte, error) {
	return ioutil.ReadFile(f.Path)
}

func (f *File) Write(data []byte) error {
	return ioutil.WriteFile(f.Path, data, 0666)
}

func (f *File) Remove() error {
	return os.Remove(f.Path)
}

func (f *File) Rename(newFile string) error {
	return os.Rename(f.Path, newFile)
}
