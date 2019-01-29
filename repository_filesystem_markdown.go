package vika

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// FilesystemMarkdownIssuesRepository stores, retrieves and otherwise manipulates issues on the filesystem.
type FilesystemMarkdownIssuesRepository struct {
	Fs Filesystem
}

// GetIssues retrieves all files from the `./.issues` folder and treats them as YAML issue definitions.
func (r FilesystemMarkdownIssuesRepository) GetIssues() ([]Issue, error) {
	files, err := r.Fs.ReadDir("./.issues")
	if err != nil {
		return nil, err
	}

	issues := make([]Issue, len(files))
	for index, file := range files {
		data, err := r.Fs.ReadFile("./.issues/" + file.Name())
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
func (r FilesystemMarkdownIssuesRepository) SaveIssue(issue *Issue) error {
	if issue.ID == "" {
		return errors.New("issue ID is empty")
	}

	// issue.Description and comments are normalized to unix newlines to force literal style marshalling in json, see https://github.com/go-yaml/yaml/issues/197
	issue.Description = string(normalizeNewlines([]byte(issue.Description)))
	for key, comment := range issue.Comments {
		comment.Message = string(normalizeNewlines([]byte(comment.Message)))
		issue.Comments[key] = comment
	}

	data, err := yaml.Marshal(&issue)
	if err != nil {
		return err
	}

	err = r.Fs.WriteFile("./.issues/"+string(issue.ID)+".yml", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetIssue finds and returns an issue from the filesystem matching the given ID.
func (r FilesystemMarkdownIssuesRepository) GetIssue(issueID ID) (Issue, error) {
	fileName := string(issueID) + ".yml"
	issue := Issue{
		ID: ID(strings.TrimSuffix(fileName, filepath.Ext(fileName))),
	}

	filePath := "./.issues/" + fileName
	data, err := r.Fs.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return issue, fmt.Errorf("issue '%s' does not exist", string(issueID))
		}
		return issue, err
	}

	err = yaml.Unmarshal([]byte(data), &issue)
	if err != nil {
		return issue, err
	}

	return issue, nil
}

// DeleteIssue deletes an issue from the filesystem matching the given ID.
func (r FilesystemMarkdownIssuesRepository) DeleteIssue(issueID ID) error {
	fileName := string(issueID) + ".yml"
	filePath := "./.issues/" + fileName

	err := r.Fs.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}
