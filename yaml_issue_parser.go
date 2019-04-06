package vika

import (
	"gopkg.in/yaml.v2"
)

func parseYmlIssue(issue *Issue, content []byte) error {
	return yaml.Unmarshal(content, issue)
}
