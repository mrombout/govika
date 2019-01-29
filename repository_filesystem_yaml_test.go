package vika

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestYamlGetIssuesWhenNoIssuesInFolderThenReturnsEmptyArray(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemYamlIssuesRepository()
	inMemFs.Mkdir(".issues", 0644)

	// act
	issues, err := repository.GetIssues()

	// assert
	assert.Empty(t, issues)
	assert.NoError(t, err)
}

func TestYamlGetIssuesWhenNoIssuesFolderThenReturnsError(t *testing.T) {
	// arrange
	repository, _ := newInMemFilesystemYamlIssuesRepository()

	// act
	issues, err := repository.GetIssues()

	// assert
	assert.Nil(t, issues)
	if assert.Error(t, err) {
		assert.EqualError(t, err, "open .issues: file does not exist")
	}
}

func TestYamlGetIssuesWhenMultipleIssuesFolderReturnsArrayWithAllIssues(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemYamlIssuesRepository()

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

func TestYamlGetIssuesWhenCantReadFileThenItReturnsAnError(t *testing.T) {
	// arrange
	expectedErr := errors.New("expected read error")
	repository := FilesystemYamlIssuesRepository{
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

func TestYamlGetIssuesWhenUnableToUnmarshalThenItReturnsAnError(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemYamlIssuesRepository()

	issueData := `invalid yaml`
	afero.WriteFile(inMemFs, ".issues/issue1.yml", []byte(issueData), 0644)

	// act
	_, err := repository.GetIssues()

	// assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml: unmarshal errors:")
}

func TestYamlSaveIssueWhenIDIsSetThenItSavesIssueToAFile(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemYamlIssuesRepository()

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

func TestYamlSaveIssueWhenIDISNotSetThenItReturnsAnError(t *testing.T) {
	// arrange
	repository, _ := newInMemFilesystemYamlIssuesRepository()

	issue := Issue{}

	// act
	err := repository.SaveIssue(&issue)

	// assert
	assert.Error(t, err)
	assert.EqualError(t, err, "issue ID is empty")
}

func TestYamlSaveIssueIDIsNotSavedToFile(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemYamlIssuesRepository()

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

	fileContent, err := afero.ReadFile(inMemFs, issueFilepath)
	assert.NoError(t, err)
	assert.NotContains(t, string(fileContent), "id: test-issue")
}

func TestYamlSaveIssueWhenCantWriteFileThenItReturnsAnError(t *testing.T) {
	// arrange
	expectedErr := errors.New("expected write error")
	repository := FilesystemYamlIssuesRepository{
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

func TestYamlGetIssueWhenIssueExistsThenItReturnsThatIssue(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemYamlIssuesRepository()

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

func TestYamlGetIssueWhenIssueDoesNotExistsThenItReturnsAnError(t *testing.T) {
	// arrange
	repository, _ := newInMemFilesystemYamlIssuesRepository()

	issueID := "issue1"

	// act
	_, err := repository.GetIssue(ID(issueID))

	// assert
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("issue '%s' does not exist", issueID))
}

func TestYamlGetIssueWhenIssueCantBeReadThenItReturnsAnError(t *testing.T) {
	// arrange
	expectedErr := errors.New("expected read error")
	repository := FilesystemYamlIssuesRepository{
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

func TestYamlGetIssueWhenUnableToUnmarshalThenItReturnsAnError(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemYamlIssuesRepository()

	issueID := "issue1"
	issueData := `invalid yaml`
	afero.WriteFile(inMemFs, fmt.Sprintf(".issues/%s.yml", issueID), []byte(issueData), 0644)

	// act
	_, err := repository.GetIssue(ID(issueID))

	// assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml: unmarshal errors:")
}

func TestYamlDeleteIssueWhenIssueExistsThenItRemovesThatIssue(t *testing.T) {
	// arrange
	repository, inMemFs := newInMemFilesystemYamlIssuesRepository()

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

func TestYamlDeleteIssueWhenIssueCantBeRemovedThenItReturnsAnError(t *testing.T) {
	// arrange
	expectedErr := errors.New("expected remove error")
	repository := FilesystemYamlIssuesRepository{
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

func newInMemFilesystemYamlIssuesRepository() (FilesystemYamlIssuesRepository, afero.Fs) {
	inMemFs := afero.NewMemMapFs()
	repository := FilesystemYamlIssuesRepository{
		Fs: AferoFilesystem{
			Fs: inMemFs,
		},
	}

	return repository, inMemFs
}
