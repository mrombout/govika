package vika

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePreamble(t *testing.T) {
	issue := Issue{}

	title := "my title"
	tokenStack := tokenStack{
		tokens: []token{
			preambleDelimiter{},
			unknownContent{Content: fmt.Sprintf("title: %s", title)},
			preambleDelimiter{},
		},
	}

	err := parsePreamble(&tokenStack, &issue)

	assert.Equal(t, title, issue.Title)
	assert.NoError(t, err)
}

func TestParsePreambleMissingFirstPreambleReturnsError(t *testing.T) {
	testCases := []struct {
		name   string
		tokens []token
	}{
		{name: "missing start preamble delimiter", tokens: []token{unknownContent{Content: "key: value"}, preambleDelimiter{}}},
		{name: "invalid preamble content", tokens: []token{preambleDelimiter{}, header1Title{Content: "key: value"}, preambleDelimiter{}}},
		{name: "missing end preamble delimiter", tokens: []token{preambleDelimiter{}, unknownContent{Content: "key: value"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			issue := Issue{}

			tokenStack := tokenStack{
				tokens: testCase.tokens,
			}

			err := parsePreamble(&tokenStack, &issue)

			assert.Error(t, err, "unexpected token")
		})
	}
}

func TestParsePreambleContainsInvalidYamlReturnsError(t *testing.T) {
	issue := Issue{}

	tokenStack := tokenStack{
		tokens: []token{
			preambleDelimiter{},
			unknownContent{Content: "not a valid yaml string"},
			preambleDelimiter{},
		},
	}

	err := parsePreamble(&tokenStack, &issue)

	assert.Error(t, err)
}

func TestParseIssueTitle(t *testing.T) {
	issue := Issue{}

	title := "my title"
	tokenStack := tokenStack{
		tokens: []token{
			header1Title{Content: title},
			emptyLine{},
		},
	}

	parseIssueTitle(&tokenStack, &issue)

	assert.Equal(t, title, issue.Title)
}

func TestParseIssueTitleUnexpectedOrMissingTokenReturnsError(t *testing.T) {
	testCases := []struct {
		name   string
		tokens []token
	}{
		{name: "missing header 1 title", tokens: []token{emptyLine{}}},
		{name: "missing empty line", tokens: []token{header1Title{}, unknownContent{Content: "issue description"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			issue := Issue{}

			tokenStack := tokenStack{
				tokens: testCase.tokens,
			}

			err := parseIssueTitle(&tokenStack, &issue)

			assert.Error(t, err, "unexpected token")
		})
	}
}

func TestParseIssueDescription(t *testing.T) {
	issue := Issue{}

	expectedText := `Lorum ipsum dolor sit amet.

Consectetur adipiscing elit.`
	tokenStack := tokenStack{
		tokens: []token{
			unknownContent{Content: "Lorum ipsum dolor sit amet."},
			emptyLine{},
			unknownContent{Content: "Consectetur adipiscing elit."},
			commentHeader{},
		},
	}

	err := parseIssueDescription(&tokenStack, &issue)

	assert.Equal(t, expectedText, issue.Description)
	assert.NoError(t, err)
}

func TestParseComments(t *testing.T) {
	tokenStack := tokenStack{
		tokens: []token{
			commentHeader{},
			emptyLine{},
		},
	}

	err := parseComments(&tokenStack)

	assert.NoError(t, err)
}

func TestParseCommentsUnexpectedOrMissingTokenReturnsError(t *testing.T) {
	testCases := []struct {
		name   string
		tokens []token
	}{
		{name: "missing comment header", tokens: []token{emptyLine{}, emptyLine{}}},
		{name: "misisng empty line", tokens: []token{commentHeader{}, commentContent{}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tokenStack := tokenStack{
				tokens: testCase.tokens,
			}

			err := parseComments(&tokenStack)

			assert.Error(t, err, "unexpected token")
		})
	}
}
