package vika

import (
	"os"

	"github.com/spf13/afero"
)

// AferoFilesystem interacts with the afero filesystem framework.
type AferoFilesystem struct {
	Fs afero.Fs
}

// ReadDir reads the directory named by dirname and returns a list of directory entries sorted by filename.
func (f AferoFilesystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return afero.ReadDir(f.Fs, dirname)
}

// ReadFile reads the file named by filename and returns the contents.
func (f AferoFilesystem) ReadFile(filename string) ([]byte, error) {
	return afero.ReadFile(f.Fs, filename)
}

// WriteFile writes data to a file named by filename.
func (f AferoFilesystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(f.Fs, filename, data, perm)
}

// Remove removes the named file or (empty) directory.
func (f AferoFilesystem) Remove(name string) error {
	return f.Fs.Remove(name)
}
