package vika

import (
	"os"

	"github.com/spf13/afero"
)

type AferoFilesystem struct {
	Fs afero.Fs
}

func (f AferoFilesystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return afero.ReadDir(f.Fs, dirname)
}

func (f AferoFilesystem) ReadFile(filename string) ([]byte, error) {
	return afero.ReadFile(f.Fs, filename)
}

func (f AferoFilesystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(f.Fs, filename, data, perm)
}

func (f AferoFilesystem) Remove(name string) error {
	return f.Fs.Remove(name)
}
