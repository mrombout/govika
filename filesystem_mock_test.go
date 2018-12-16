package vika

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	// arrange
	fileInfo := mockFileInfo{
		name: "test-name",
	}

	// act
	name := fileInfo.Name()

	// assert
	assert.Equal(t, fileInfo.name, name)
}

func TestSize(t *testing.T) {
	// arrange
	fileInfo := mockFileInfo{
		size: 64,
	}

	// act
	size := fileInfo.Size()

	// assert
	assert.Equal(t, fileInfo.size, size)
}

func TestMode(t *testing.T) {
	// arrange
	fileInfo := mockFileInfo{
		mode: 0644,
	}

	// act
	mode := fileInfo.Mode()

	// assert
	assert.Equal(t, fileInfo.mode, mode)
}

func TestModtime(t *testing.T) {
	// arrange
	fileInfo := mockFileInfo{
		modTime: time.Now(),
	}

	// act
	modTime := fileInfo.ModTime()

	// assert
	assert.Equal(t, fileInfo.modTime, modTime)
}

func TestIsDir(t *testing.T) {
	// arrange
	fileInfo := mockFileInfo{
		isDir: true,
	}

	// act
	isDir := fileInfo.IsDir()

	// assert
	assert.Equal(t, fileInfo.isDir, isDir)
}

func TestSys(t *testing.T) {
	// arrange
	fileInfo := mockFileInfo{
		sys: 1234,
	}

	// act
	sys := fileInfo.Sys()

	// assert
	assert.Equal(t, fileInfo.sys, sys)
}
