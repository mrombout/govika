package vika

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
