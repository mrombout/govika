package vika

import "os"

// Filesystem abstracts away any filesystem interaction.
type Filesystem interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Remove(name string) error
}
