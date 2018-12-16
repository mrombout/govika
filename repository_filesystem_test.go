package vika

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetIssuesWhenNoIssuesInFolderThenReturnsEmptyArray(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemIssuesRepository()
	inMemFs.Mkdir(".issues", 0644)

	// act
	issues, err := repository.GetIssues()

	// assert
	assert.Empty(t, issues)
	assert.NoError(t, err)
}

func TestGetIssuesWhenNoIssuesFolderThenReturnsError(t *testing.T) {
	// arrange
	repository, _ := newInMemFilesystemIssuesRepository()

	// act
	issues, err := repository.GetIssues()

	// assert
	assert.Nil(t, issues)
	if assert.Error(t, err) {
		assert.EqualError(t, err, "open .issues: file does not exist")
	}
}

func TestGetIssuesWhenMultipleIssuesFolderReturnsArrayWithAllIssues(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemIssuesRepository()

	firstIssueTitle := "Issue 1"
	firstIssueData := fmt.Sprintf(`id: issue1
title: %s
description: This is issue 1.
author: Test`, firstIssueTitle)
	afero.WriteFile(inMemFs, ".issues/issue1.yml", []byte(firstIssueData), 0644)
	secondIssueTitle := "Issue 2"
	secondIssueData := fmt.Sprintf(`id: issue2
title: %s
description: This is issue 2.
author: Test`, secondIssueTitle)
	afero.WriteFile(inMemFs, ".issues/issue2.yml", []byte(secondIssueData), 0644)

	// act
	issues, err := repository.GetIssues()

	// assert
	assert.Len(t, issues, 2)
	assert.Equal(t, firstIssueTitle, issues[0].Title)
	assert.Equal(t, secondIssueTitle, issues[1].Title)
	assert.NoError(t, err)
}

func TestGetIssuesWhenCantReadFileThenItReturnsAnError(t *testing.T) {
	// arrange
	expectedErr := errors.New("expected read error")
	repository := FilesystemIssuesRepository{
		Fs: mockFilesystem{
			readFileReturnError: expectedErr,
			readDirReturnFileInfo: []os.FileInfo{
				mockFileInfo{name: "issue1.yml"},
			},
		},
	}

	// act
	_, err := repository.GetIssues()

	// assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestGetIssuesWhenUnableToUnmarshalThenItReturnsAnError(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemIssuesRepository()

	issueData := `invalid yaml`
	afero.WriteFile(inMemFs, ".issues/issue1.yml", []byte(issueData), 0644)

	// act
	_, err := repository.GetIssues()

	// assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml: unmarshal errors:")
}

func TestSaveIssueWhenIDIsSetThenItSavesIssueToAFile(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemIssuesRepository()

	issueID := "test-issue"
	issue := Issue{
		ID: ID(issueID),
	}
	issueFilepath := fmt.Sprintf(".issues/%s.yml", issueID)

	// act
	err := repository.SaveIssue(&issue)

	// assert
	assert.NoError(t, err)
	assertFileExists(t, inMemFs, issueFilepath)
}

func TestSaveIssueWhenIDISNotSetThenItReturnsAnError(t *testing.T) {
	// arrange
	repository, _ := newInMemFilesystemIssuesRepository()

	issue := Issue{}

	// act
	err := repository.SaveIssue(&issue)

	// assert
	assert.Error(t, err)
	assert.EqualError(t, err, "issue ID is empty")
}

func TestSaveIssueWhenCantWriteFileThenItReturnsAnError(t *testing.T) {
	// arrange
	expectedErr := errors.New("expected write error")
	repository := FilesystemIssuesRepository{
		Fs: mockFilesystem{
			writeFileReturnError: expectedErr,
		},
	}

	issue := Issue{
		ID: "issue1",
	}

	// act
	err := repository.SaveIssue(&issue)

	// assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestGetIssueWhenIssueExistsThenItReturnsThatIssue(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemIssuesRepository()

	issueID := "issue1"
	issueTitle := "Issue 1"
	issueData := fmt.Sprintf(`id: %s
title: %s
description: This is issue 1.
author: Test`, issueID, issueTitle)
	afero.WriteFile(inMemFs, fmt.Sprintf(".issues/%s.yml", issueID), []byte(issueData), 0644)

	// act
	issue, err := repository.GetIssue(ID(issueID))

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, issue)
	assert.Equal(t, issueTitle, issue.Title)
}

func TestGetIssueWhenIssueDoesNotExistsThenItReturnsAnError(t *testing.T) {
	// arrange
	repository, _ := newInMemFilesystemIssuesRepository()

	issueID := "issue1"

	// act
	_, err := repository.GetIssue(ID(issueID))

	// assert
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("issue '%s' does not exist", issueID))
}

func TestGetIssueWhenIssueCantBeReadThenItReturnsAnError(t *testing.T) {
	// arrange
	expectedErr := errors.New("expected read error")
	repository := FilesystemIssuesRepository{
		Fs: mockFilesystem{
			readFileReturnError: expectedErr,
		},
	}

	// act
	_, err := repository.GetIssue("issue1")

	// assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestGetIssueWhenUnableToUnmarshalThenItReturnsAnError(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemIssuesRepository()

	issueID := "issue1"
	issueData := `invalid yaml`
	afero.WriteFile(inMemFs, fmt.Sprintf(".issues/%s.yml", issueID), []byte(issueData), 0644)

	// act
	_, err := repository.GetIssue(ID(issueID))

	// assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml: unmarshal errors:")
}

func TestDeleteIssueWhenIssueExistsThenItRemovesThatIssue(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemIssuesRepository()

	issueID := "issue1"
	issueData := fmt.Sprintf(`id: %s
	title: Issue 1
	description: This is issue 1.
	author: Test`, issueID)
	issueFilepath := fmt.Sprintf(".issues/%s.yml", issueID)
	afero.WriteFile(inMemFs, issueFilepath, []byte(issueData), 0644)

	// act
	err := repository.DeleteIssue(ID(issueID))

	// assert
	assert.NoError(t, err)
	assertFileNotExists(t, inMemFs, issueFilepath)
}

func TestDeleteIssueWhenIssueCantBeRemovedThenItReturnsAnError(t *testing.T) {
	// arrange
	expectedErr := errors.New("expected remove error")
	repository := FilesystemIssuesRepository{
		Fs: mockFilesystem{
			removeReturnError: expectedErr,
		},
	}

	// act
	err := repository.DeleteIssue("issue1")

	// assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
}

func newInMemFilesystemIssuesRepository() (FilesystemIssuesRepository, afero.Fs) {
	inMemFs := afero.NewMemMapFs()
	repository := FilesystemIssuesRepository{
		Fs: AferoFilesystem{
			Fs: inMemFs,
		},
	}

	return repository, inMemFs
}

func assertFileExists(t *testing.T, fs afero.Fs, filepath string) {
	t.Helper()

	_, err := fs.Stat(filepath)
	assert.NoError(t, err)
}

func assertFileNotExists(t *testing.T, fs afero.Fs, filepath string) {
	t.Helper()

	_, err := fs.Stat(filepath)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}
