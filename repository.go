package vika

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Comment represents a comment on an issue.
type Comment struct {
	Author  string
	Message string
}

// Label represents a label associated with one or more issues.
type Label string

// ID represents an issue ID. It may not contain any spaces or special characters that are also illegal in file paths.
type ID string

// Issue represents something that needs to be done or discussed in a project (e.g. a bug, a task, an issue, a feature
// request, etc)
type Issue struct {
	ID          ID
	Title       string
	Description string
	Author      string
	Milestone   string
	Comments    []Comment
	Labels      []Label
}

// IssuesRepository stores, retrieves and otherwise manipulates issues of a project.
type IssuesRepository interface {
	GetIssues() ([]Issue, error)
	SaveIssue(issue *Issue) error
	GetIssue(ID ID) (Issue, error)
	DeleteIssue(ID ID) error
}

// FilesystemIssuesRepository stores, retrieves and otherwise manipulates issues on the filesystem.
type FilesystemIssuesRepository struct {
}

// GetIssues retrieves all files from the `./.issues` folder and treats them as YAML issue definitions.
func (FilesystemIssuesRepository) GetIssues() ([]Issue, error) {
	files, err := ioutil.ReadDir("./.issues")
	if err != nil {
		return nil, err
	}

	issues := make([]Issue, len(files))
	for index, file := range files {
		data, err := ioutil.ReadFile("./.issues/" + file.Name())
		if err != nil {
			return nil, err
		}

		issue := Issue{
			ID: ID(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))),
		}
		err = yaml.Unmarshal([]byte(data), &issue)
		if err != nil {
			return nil, err
		}
		issues[index] = issue
	}

	return issues, nil
}

// SaveIssue saves an issue to the `./.issues` folder. The existing issue will be overwritten.
func (FilesystemIssuesRepository) SaveIssue(issue *Issue) error {
	data, err := yaml.Marshal(&issue)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./.issues/"+string(issue.ID)+".yml", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetIssue finds and returns an issue from the filesystem matching the given ID.
func (FilesystemIssuesRepository) GetIssue(issueID ID) (Issue, error) {
	fileName := string(issueID) + ".yml"
	issue := Issue{
		ID: ID(strings.TrimSuffix(fileName, filepath.Ext(fileName))),
	}

	filePath := "./.issues/" + fileName
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return issue, err
	}

	err = yaml.Unmarshal([]byte(data), &issue)
	if err != nil {
		return issue, err
	}

	return issue, nil
}

// DeleteIssue deletes an issue from the filesystem matching the given ID.
func (FilesystemIssuesRepository) DeleteIssue(issueID ID) error {
	fileName := string(issueID) + ".yml"
	filePath := "./.issues/" + fileName

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}
