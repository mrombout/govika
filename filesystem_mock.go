package vika

import (
	"os"
	"time"
)

type mockFilesystem struct {
	readDirReturnFileInfo []os.FileInfo
	readFileReturnError   error
	writeFileReturnError  error
	removeReturnError     error
}

func (f mockFilesystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return f.readDirReturnFileInfo, nil
}

func (f mockFilesystem) ReadFile(filename string) ([]byte, error) {
	return []byte{}, f.readFileReturnError
}

func (f mockFilesystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return f.writeFileReturnError
}

func (f mockFilesystem) Remove(name string) error {
	return f.removeReturnError
}

type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (m mockFileInfo) Name() string {
	return m.name
}

func (m mockFileInfo) Size() int64 {
	return m.size
}

func (m mockFileInfo) Mode() os.FileMode {
	return m.mode
}

func (m mockFileInfo) ModTime() time.Time {
	return m.modTime
}

func (m mockFileInfo) IsDir() bool {
	return m.isDir
}

func (m mockFileInfo) Sys() interface{} {
	return m.sys
}
