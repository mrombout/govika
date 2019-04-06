package vika

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseYmlIssueParsesTitle(t *testing.T) {
	issue := Issue{}

	title := "Issue 1"

	content := fmt.Sprintf(`id: issue1
title: %s
description: This is issue 1.
author: Test`, title)

	parseYmlIssue(&issue, []byte(content))

	assert.Equal(t, title, issue.Title)
}
